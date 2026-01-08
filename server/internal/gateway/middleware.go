package gateway

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
)

type AuthProvider interface {
	VerifyToken(ctx context.Context, token string) (*TokenInfo, error)
	GetTenantInfo(ctx context.Context, tenantID string) (*TenantInfo, error)
	GenerateInternalToken(ctx context.Context, info *TokenInfo) (string, error)
}

type TokenInfo struct {
	UserID      string
	TenantID    string
	Permissions []string
}

type TenantInfo struct {
	ID             string
	DefaultService string
	IsActive       bool
}

func AuthMiddleware(authProvider AuthProvider, cache *Cache, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		opaqueToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Check cache for token info
		var tokenInfo *TokenInfo
		if cached, ok := cache.Get("token:" + opaqueToken); ok {
			tokenInfo = cached.(*TokenInfo)
		} else {
			var err error
			tokenInfo, err = authProvider.VerifyToken(c.Request.Context(), opaqueToken)
			if err != nil {
				log.Error("Failed to verify token", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}
			cache.Set("token:"+opaqueToken, tokenInfo, 5*time.Minute)
		}

		// Verify tenant
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = tokenInfo.TenantID
		}

		// Check tenant in cache
		var tenantInfo *TenantInfo
		if cached, ok := cache.Get("tenant:" + tenantID); ok {
			tenantInfo = cached.(*TenantInfo)
		} else {
			var err error
			tenantInfo, err = authProvider.GetTenantInfo(c.Request.Context(), tenantID)
			if err != nil {
				log.Error("Failed to get tenant info", zap.Error(err))
				// Failover logic mentioned in point 7 will be handled in routing
			} else {
				cache.Set("tenant:"+tenantID, tenantInfo, 10*time.Minute)
			}
		}

		// Generate internal token
		internalToken, err := authProvider.GenerateInternalToken(c.Request.Context(), tokenInfo)
		if err != nil {
			log.Error("Failed to generate internal token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}

		// Inject headers
		c.Request.Header.Set("X-Tenant-ID", tenantID)
		c.Request.Header.Set("X-Internal-Token", internalToken)
		c.Set("tenant_info", tenantInfo) // Pass tenant info to next middleware

		c.Next()
	}
}
