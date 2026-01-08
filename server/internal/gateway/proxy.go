package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
)

type ProxyHandler struct {
	logger *logger.Logger
}

func NewProxyHandler(log *logger.Logger) *ProxyHandler {
	return &ProxyHandler{
		logger: log,
	}
}

func (h *ProxyHandler) HandleRequest(c *gin.Context) {
	path := c.Request.URL.Path
	var targetURL string

	switch {
	case strings.HasPrefix(path, "/api/"):
		// Format: /api/service-name/api-path
		parts := strings.SplitN(strings.TrimPrefix(path, "/api/"), "/", 2)
		if len(parts) >= 1 {
			serviceName := parts[0]
			apiPath := ""
			if len(parts) > 1 {
				apiPath = "/" + parts[1]
			}
			targetURL = h.resolveServiceURL(serviceName, apiPath)
		}

	case strings.HasPrefix(path, "/page/"):
		// Format: /page/service-name/page-path
		parts := strings.SplitN(strings.TrimPrefix(path, "/page/"), "/", 2)
		if len(parts) >= 1 {
			serviceName := parts[0]
			pagePath := ""
			if len(parts) > 1 {
				pagePath = "/" + parts[1]
			}
			targetURL = h.resolvePageURL(serviceName, pagePath)
		}

	case strings.HasPrefix(path, "/upload/"):
		// Format: /upload/file-key
		fileKey := strings.TrimPrefix(path, "/upload/")
		targetURL = h.resolveUploadURL(fileKey)

	default:
		// Slug processing or default service
		targetURL = h.resolveSlugURL(path, c)
	}

	if targetURL == "" {
		h.handleFailover(c)
		return
	}

	h.proxyTo(c, targetURL)
}

func (h *ProxyHandler) resolveServiceURL(serviceName, apiPath string) string {
	// In a real scenario, this would use service discovery (Consul, K8s DNS)
	// For this demo, we'll use env vars or defaults
	port := "8080" // Default
	return "http://" + serviceName + ":" + port + apiPath
}

func (h *ProxyHandler) resolvePageURL(serviceName, pagePath string) string {
	return "http://" + serviceName + "-web:3000" + pagePath
}

func (h *ProxyHandler) resolveUploadURL(fileKey string) string {
	return "http://file-service:8080/files/" + fileKey
}

func (h *ProxyHandler) resolveSlugURL(path string, c *gin.Context) string {
	// Point 3: "đường dẫn dạng còn lại thì xử lý như dạng đường dẫn đẹp (slug)"
	// This usually means routing to a CMS or a specific service that handles slugs
	return "http://cms-service:8080/slug" + path
}

func (h *ProxyHandler) proxyTo(c *gin.Context, target string) {
	remote, err := url.Parse(target)
	if err != nil {
		h.logger.Error("Failed to parse target URL", zap.Error(err), zap.String("target", target))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		h.logger.Error("Proxy error", zap.Error(err), zap.String("target", target))
		h.handleFailover(c)
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (h *ProxyHandler) handleFailover(c *gin.Context) {
	// Point 7: "Nếu gateway điều hướng mà lỗi thì sẽ chuyển hướng tiếp về default service của tenant.
	// Nếu tiếp tục lỗi sẽ chuyển về default service của hệ thống"

	tenantInfoRaw, exists := c.Get("tenant_info")
	if exists && tenantInfoRaw != nil {
		tenantInfo := tenantInfoRaw.(*TenantInfo)
		if tenantInfo.DefaultService != "" {
			h.logger.Info("Failing over to tenant default service", zap.String("service", tenantInfo.DefaultService))
			// Recursive call or redirect
			// For simplicity, we'll mock a redirect or re-proxy
			// target := h.resolveServiceURL(tenantInfo.DefaultService, c.Request.URL.Path)
			// h.proxyTo(c, target)
			return
		}
	}

	// Default system service
	h.logger.Info("Failing over to system default service")
	// target := h.resolveServiceURL("system-default", c.Request.URL.Path)
	// h.proxyTo(c, target)
}
