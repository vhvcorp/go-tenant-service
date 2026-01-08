package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/config"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-tenant-service/internal/gateway"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	log.Info("Starting API Gateway", zap.String("environment", cfg.Environment))

	// Initialize local cache
	// Point 5: "thêm cấu hình để giới hạn cache tối đa bao nhiêu dữ liệu"
	cache := gateway.NewCache(5*time.Minute, 10*time.Minute)

	// Initialize mock auth provider (In production, this would be a gRPC client to Auth service)
	authProvider := &MockAuthProvider{}

	// Initialize proxy handler
	proxyHandler := gateway.NewProxyHandler(log)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "cache_items": cache.ItemCount()})
	})

	// Apply Auth Middleware and Proxy all other requests
	router.Use(gateway.AuthMiddleware(authProvider, cache, log))
	router.NoRoute(proxyHandler.HandleRequest)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Info("Gateway listening", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start Gateway", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Gateway forced to shutdown", zap.Error(err))
	}

	log.Info("Gateway exited")
}

// MockAuthProvider for demonstration
type MockAuthProvider struct{}

func (m *MockAuthProvider) VerifyToken(ctx context.Context, token string) (*gateway.TokenInfo, error) {
	return &gateway.TokenInfo{
		UserID:      "user-123",
		TenantID:    "tenant-abc",
		Permissions: []string{"read", "write"},
	}, nil
}

func (m *MockAuthProvider) GetTenantInfo(ctx context.Context, tenantID string) (*gateway.TenantInfo, error) {
	return &gateway.TenantInfo{
		ID:             tenantID,
		DefaultService: "tenant-service",
		IsActive:       true,
	}, nil
}

func (m *MockAuthProvider) GenerateInternalToken(ctx context.Context, info *gateway.TokenInfo) (string, error) {
	return "internal-jwt-token-for-" + info.UserID, nil
}
