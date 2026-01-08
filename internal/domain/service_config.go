package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServiceConfig defines service endpoint configuration for a tenant
type ServiceConfig struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID    string             `bson:"tenantId" json:"tenant_id"`
	ServiceName string             `bson:"serviceName" json:"service_name"` // e.g., "user", "auth", "notification"

	// Primary endpoint
	PrimaryEndpoint ServiceEndpoint `bson:"primaryEndpoint" json:"primary_endpoint"`

	// Fallback chain (tried in order if primary fails)
	FallbackChain []ServiceEndpoint `bson:"fallbackChain,omitempty" json:"fallback_chain,omitempty"`

	// Default service URL (used if no tenant-specific config)
	DefaultServiceURL string `bson:"defaultServiceUrl,omitempty" json:"default_service_url,omitempty"`

	// Health check configuration
	HealthCheck HealthCheckConfig `bson:"healthCheck,omitempty" json:"health_check,omitempty"`

	// Load balancing strategy
	LoadBalanceStrategy string `bson:"loadBalanceStrategy,omitempty" json:"load_balance_strategy,omitempty"` // "round-robin", "random", "weighted"

	// Metadata
	IsActive  bool              `bson:"isActive" json:"is_active"`
	Metadata  map[string]string `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time         `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time         `bson:"updatedAt" json:"updated_at"`
}

// ServiceEndpoint represents a single service endpoint
type ServiceEndpoint struct {
	URL      string            `bson:"url" json:"url"`                             // e.g., "http://user-service:8080"
	Priority int               `bson:"priority" json:"priority"`                   // Lower = higher priority (0 = primary)
	Weight   int               `bson:"weight,omitempty" json:"weight,omitempty"`   // For weighted load balancing
	Timeout  int               `bson:"timeout,omitempty" json:"timeout,omitempty"` // Timeout in seconds
	Headers  map[string]string `bson:"headers,omitempty" json:"headers,omitempty"` // Custom headers
	IsActive bool              `bson:"isActive" json:"is_active"`
}

// HealthCheckConfig defines how to check service health
type HealthCheckConfig struct {
	Enabled       bool   `bson:"enabled" json:"enabled"`
	Path          string `bson:"path,omitempty" json:"path,omitempty"`                    // e.g., "/health"
	Method        string `bson:"method,omitempty" json:"method,omitempty"`                // "GET", "POST"
	Interval      int    `bson:"interval,omitempty" json:"interval,omitempty"`            // Check interval in seconds
	Timeout       int    `bson:"timeout,omitempty" json:"timeout,omitempty"`              // Health check timeout
	FailThreshold int    `bson:"failThreshold,omitempty" json:"fail_threshold,omitempty"` // Failures before marking unhealthy
}

// DefaultServiceConfig holds system-wide default configurations
type DefaultServiceConfig struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ServiceName string             `bson:"serviceName" json:"service_name"`
	ServiceURL  string             `bson:"serviceUrl" json:"service_url"`

	// Default health check
	HealthCheck HealthCheckConfig `bson:"healthCheck,omitempty" json:"health_check,omitempty"`

	// Fallback behavior
	FallbackToDefault bool `bson:"fallbackToDefault" json:"fallback_to_default"`

	CreatedAt time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updated_at"`
}

// ServiceStatus represents real-time service health status
type ServiceStatus struct {
	ServiceName      string    `json:"service_name"`
	EndpointURL      string    `json:"endpoint_url"`
	IsHealthy        bool      `json:"is_healthy"`
	LastCheckTime    time.Time `json:"last_check_time"`
	ConsecutiveFails int       `json:"consecutive_fails"`
	ResponseTime     int64     `json:"response_time_ms"`
	ErrorMessage     string    `json:"error_message,omitempty"`
}

// ServiceRegistryEntry combines config and status
type ServiceRegistryEntry struct {
	Config ServiceConfig `json:"config"`
	Status ServiceStatus `json:"status"`
}

// FallbackChainResult contains the result of fallback chain resolution
type FallbackChainResult struct {
	SelectedEndpoint ServiceEndpoint `json:"selected_endpoint"`
	TriedEndpoints   []string        `json:"tried_endpoints"`
	FallbackLevel    int             `json:"fallback_level"` // 0 = primary, 1+ = fallback
	Success          bool            `json:"success"`
	Error            string          `json:"error,omitempty"`
}

// ServiceDiscoveryRequest for dynamic service discovery
type ServiceDiscoveryRequest struct {
	TenantID    string `json:"tenant_id"`
	ServiceName string `json:"service_name"`
	Version     string `json:"version,omitempty"` // Optional version requirement
}

// ServiceDiscoveryResponse contains discovered service info
type ServiceDiscoveryResponse struct {
	ServiceName string            `json:"service_name"`
	Endpoint    ServiceEndpoint   `json:"endpoint"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Constants for load balancing strategies
const (
	LoadBalanceRoundRobin = "round-robin"
	LoadBalanceRandom     = "random"
	LoadBalanceWeighted   = "weighted"
	LoadBalanceLeastConn  = "least-conn"
)

// Constants for service names
const (
	ServiceAuth         = "auth"
	ServiceUser         = "user"
	ServiceTenant       = "tenant"
	ServiceNotification = "notification"
	ServiceCMS          = "cms"
	ServiceConfig       = "config"
)

// Default health check values
const (
	DefaultHealthCheckPath     = "/health"
	DefaultHealthCheckMethod   = "GET"
	DefaultHealthCheckInterval = 30 // seconds
	DefaultHealthCheckTimeout  = 5  // seconds
	DefaultFailThreshold       = 3  // consecutive failures
)

// Validate validates the ServiceConfig
func (sc *ServiceConfig) Validate() error {
	if sc.TenantID == "" {
		return ErrTenantIDRequired
	}
	if sc.ServiceName == "" {
		return ErrServiceNameRequired
	}
	if sc.PrimaryEndpoint.URL == "" {
		return ErrPrimaryEndpointRequired
	}
	return nil
}

// GetActiveEndpoints returns all active endpoints (primary + fallbacks)
func (sc *ServiceConfig) GetActiveEndpoints() []ServiceEndpoint {
	endpoints := []ServiceEndpoint{}

	if sc.PrimaryEndpoint.IsActive {
		endpoints = append(endpoints, sc.PrimaryEndpoint)
	}

	for _, endpoint := range sc.FallbackChain {
		if endpoint.IsActive {
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}

// GetEndpointByPriority returns endpoints sorted by priority
func (sc *ServiceConfig) GetEndpointByPriority() []ServiceEndpoint {
	endpoints := sc.GetActiveEndpoints()

	// Sort by priority (lower number = higher priority)
	for i := 0; i < len(endpoints); i++ {
		for j := i + 1; j < len(endpoints); j++ {
			if endpoints[j].Priority < endpoints[i].Priority {
				endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
			}
		}
	}

	return endpoints
}

// IsHealthCheckEnabled checks if health checking is enabled
func (sc *ServiceConfig) IsHealthCheckEnabled() bool {
	return sc.HealthCheck.Enabled
}

// GetHealthCheckInterval returns health check interval with default
func (sc *ServiceConfig) GetHealthCheckInterval() int {
	if sc.HealthCheck.Interval <= 0 {
		return DefaultHealthCheckInterval
	}
	return sc.HealthCheck.Interval
}

// Custom errors
var (
	ErrTenantIDRequired        = NewValidationError("tenant_id is required")
	ErrServiceNameRequired     = NewValidationError("service_name is required")
	ErrPrimaryEndpointRequired = NewValidationError("primary_endpoint is required")
	ErrServiceNotFound         = NewNotFoundError("service configuration not found")
	ErrNoHealthyEndpoint       = NewServiceError("no healthy endpoint available")
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Message: msg}
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Message: msg}
}

// ServiceError represents a service error
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewServiceError(msg string) *ServiceError {
	return &ServiceError{Message: msg}
}
