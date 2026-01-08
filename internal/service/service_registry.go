package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"github.com/vhvplatform/go-tenant-service/internal/repository"
)

// ServiceRegistry manages service discovery and routing
type ServiceRegistry struct {
	repo           *repository.ServiceConfigRepository
	healthStatus   map[string]*domain.ServiceStatus // key: tenantID:serviceName:url
	statusMutex    sync.RWMutex
	loadBalanceIdx map[string]int // key: tenantID:serviceName (for round-robin)
	lbMutex        sync.Mutex
	logger         logger.Logger
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(repo *repository.ServiceConfigRepository, log logger.Logger) *ServiceRegistry {
	return &ServiceRegistry{
		repo:           repo,
		healthStatus:   make(map[string]*domain.ServiceStatus),
		statusMutex:    sync.RWMutex{},
		loadBalanceIdx: make(map[string]int),
		lbMutex:        sync.Mutex{},
		logger:         log,
	}
}

// GetServiceURL resolves the best service URL for a tenant and service
// It follows the fallback chain: tenant config -> default config -> error
func (s *ServiceRegistry) GetServiceURL(ctx context.Context, tenantID, serviceName string) (*domain.FallbackChainResult, error) {
	result := &domain.FallbackChainResult{
		TenantID:    tenantID,
		ServiceName: serviceName,
		AttemptedAt: time.Now(),
	}

	// Try tenant-specific configuration
	config, err := s.repo.FindByTenantAndService(ctx, tenantID, serviceName)
	if err != nil {
		s.logger.Error("Failed to find tenant service config",
			"tenantId", tenantID,
			"service", serviceName,
			"error", err)
	}

	if config != nil && config.IsActive {
		url, endpoint := s.selectEndpoint(config)
		if url != "" {
			result.ResolvedURL = url
			result.UsedEndpoint = endpoint
			result.IsDefault = false
			result.Success = true
			return result, nil
		}
	}

	// Fallback to default configuration
	defaultConfig, err := s.repo.GetDefaultConfig(ctx, serviceName)
	if err != nil {
		s.logger.Error("Failed to find default service config",
			"service", serviceName,
			"error", err)
		result.Success = false
		result.Error = fmt.Sprintf("failed to resolve service URL: %v", err)
		return result, domain.ErrServiceNotFound
	}

	if defaultConfig == nil {
		result.Success = false
		result.Error = "no configuration found for service"
		return result, domain.ErrServiceNotFound
	}

	result.ResolvedURL = defaultConfig.DefaultURL
	result.IsDefault = true
	result.Success = true
	return result, nil
}

// selectEndpoint selects the best endpoint based on load balancing strategy
func (s *ServiceRegistry) selectEndpoint(config *domain.ServiceConfig) (string, *domain.ServiceEndpoint) {
	endpoints := config.GetActiveEndpoints()
	if len(endpoints) == 0 {
		// No active endpoints, try primary even if inactive
		if config.PrimaryEndpoint.URL != "" {
			return config.PrimaryEndpoint.URL, &config.PrimaryEndpoint
		}
		// Last resort: use default service URL
		return config.DefaultServiceURL, nil
	}

	switch config.LoadBalanceStrategy {
	case domain.LoadBalanceRoundRobin:
		return s.roundRobinSelect(config.TenantID, config.ServiceName, endpoints)
	case domain.LoadBalanceRandom:
		return s.randomSelect(endpoints)
	case domain.LoadBalanceWeighted:
		return s.weightedSelect(endpoints)
	case domain.LoadBalanceLeastConn:
		// For now, fallback to round-robin (least-conn requires connection tracking)
		return s.roundRobinSelect(config.TenantID, config.ServiceName, endpoints)
	default:
		// Default to round-robin
		return s.roundRobinSelect(config.TenantID, config.ServiceName, endpoints)
	}
}

// roundRobinSelect selects endpoint using round-robin algorithm
func (s *ServiceRegistry) roundRobinSelect(tenantID, serviceName string, endpoints []*domain.ServiceEndpoint) (string, *domain.ServiceEndpoint) {
	if len(endpoints) == 0 {
		return "", nil
	}

	key := fmt.Sprintf("%s:%s", tenantID, serviceName)

	s.lbMutex.Lock()
	defer s.lbMutex.Unlock()

	idx := s.loadBalanceIdx[key]
	endpoint := endpoints[idx%len(endpoints)]

	s.loadBalanceIdx[key] = (idx + 1) % len(endpoints)

	return endpoint.URL, endpoint
}

// randomSelect selects a random active endpoint
func (s *ServiceRegistry) randomSelect(endpoints []*domain.ServiceEndpoint) (string, *domain.ServiceEndpoint) {
	if len(endpoints) == 0 {
		return "", nil
	}

	idx := rand.Intn(len(endpoints))
	endpoint := endpoints[idx]
	return endpoint.URL, endpoint
}

// weightedSelect selects endpoint based on weight
func (s *ServiceRegistry) weightedSelect(endpoints []*domain.ServiceEndpoint) (string, *domain.ServiceEndpoint) {
	if len(endpoints) == 0 {
		return "", nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, ep := range endpoints {
		totalWeight += ep.Weight
	}

	if totalWeight == 0 {
		// If no weights, fallback to random
		return s.randomSelect(endpoints)
	}

	// Select based on weight
	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, ep := range endpoints {
		cumulative += ep.Weight
		if r < cumulative {
			return ep.URL, ep
		}
	}

	// Fallback (should not reach here)
	return endpoints[0].URL, endpoints[0]
}

// ResolveFallbackChain attempts to resolve a working endpoint through the fallback chain
func (s *ServiceRegistry) ResolveFallbackChain(ctx context.Context, config *domain.ServiceConfig) (*domain.FallbackChainResult, error) {
	result := &domain.FallbackChainResult{
		TenantID:      config.TenantID,
		ServiceName:   config.ServiceName,
		AttemptedURLs: []string{},
		AttemptedAt:   time.Now(),
	}

	// Try primary endpoint
	if config.PrimaryEndpoint.IsActive && s.isEndpointHealthy(config.TenantID, config.ServiceName, config.PrimaryEndpoint.URL) {
		result.ResolvedURL = config.PrimaryEndpoint.URL
		result.UsedEndpoint = &config.PrimaryEndpoint
		result.IsDefault = false
		result.Success = true
		return result, nil
	}
	result.AttemptedURLs = append(result.AttemptedURLs, config.PrimaryEndpoint.URL)

	// Try fallback chain
	for _, endpoint := range config.FallbackChain {
		if endpoint.IsActive && s.isEndpointHealthy(config.TenantID, config.ServiceName, endpoint.URL) {
			result.ResolvedURL = endpoint.URL
			result.UsedEndpoint = &endpoint
			result.IsDefault = false
			result.Success = true
			return result, nil
		}
		result.AttemptedURLs = append(result.AttemptedURLs, endpoint.URL)
	}

	// Last resort: use default service URL
	if config.DefaultServiceURL != "" {
		result.ResolvedURL = config.DefaultServiceURL
		result.IsDefault = true
		result.Success = true
		return result, nil
	}

	result.Success = false
	result.Error = "all endpoints in fallback chain are unhealthy"
	return result, fmt.Errorf("no healthy endpoints available")
}

// UpdateHealthStatus updates the health status of an endpoint
func (s *ServiceRegistry) UpdateHealthStatus(tenantID, serviceName, url string, healthy bool) {
	key := fmt.Sprintf("%s:%s:%s", tenantID, serviceName, url)

	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	status, exists := s.healthStatus[key]
	if !exists {
		status = &domain.ServiceStatus{
			TenantID:         tenantID,
			ServiceName:      serviceName,
			EndpointURL:      url,
			IsHealthy:        healthy,
			LastChecked:      time.Now(),
			ConsecutiveFails: 0,
		}
		s.healthStatus[key] = status
		return
	}

	status.LastChecked = time.Now()
	if healthy {
		status.IsHealthy = true
		status.ConsecutiveFails = 0
		status.LastSuccessful = time.Now()
	} else {
		status.ConsecutiveFails++
		status.LastFailure = time.Now()
		status.LastError = "health check failed"

		// Mark as unhealthy after threshold
		if status.ConsecutiveFails >= 3 {
			status.IsHealthy = false
		}
	}
}

// isEndpointHealthy checks if an endpoint is currently healthy
func (s *ServiceRegistry) isEndpointHealthy(tenantID, serviceName, url string) bool {
	key := fmt.Sprintf("%s:%s:%s", tenantID, serviceName, url)

	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()

	status, exists := s.healthStatus[key]
	if !exists {
		// If no status recorded, assume healthy
		return true
	}

	return status.IsHealthy
}

// GetHealthStatus returns the current health status of an endpoint
func (s *ServiceRegistry) GetHealthStatus(tenantID, serviceName, url string) *domain.ServiceStatus {
	key := fmt.Sprintf("%s:%s:%s", tenantID, serviceName, url)

	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()

	status, exists := s.healthStatus[key]
	if !exists {
		return &domain.ServiceStatus{
			TenantID:    tenantID,
			ServiceName: serviceName,
			EndpointURL: url,
			IsHealthy:   true, // Default to healthy
			LastChecked: time.Now(),
		}
	}

	// Return a copy to avoid race conditions
	statusCopy := *status
	return &statusCopy
}

// GetAllHealthStatus returns all health status records
func (s *ServiceRegistry) GetAllHealthStatus() []*domain.ServiceStatus {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()

	statuses := make([]*domain.ServiceStatus, 0, len(s.healthStatus))
	for _, status := range s.healthStatus {
		statusCopy := *status
		statuses = append(statuses, &statusCopy)
	}

	return statuses
}

// InvalidateCache clears the health status cache for a specific endpoint
func (s *ServiceRegistry) InvalidateCache(tenantID, serviceName, url string) {
	key := fmt.Sprintf("%s:%s:%s", tenantID, serviceName, url)

	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	delete(s.healthStatus, key)
}

// InvalidateAllCache clears all health status cache
func (s *ServiceRegistry) InvalidateAllCache() {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	s.healthStatus = make(map[string]*domain.ServiceStatus)
}

// CreateOrUpdateServiceConfig creates or updates a service configuration
func (s *ServiceRegistry) CreateOrUpdateServiceConfig(ctx context.Context, config *domain.ServiceConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}

	return s.repo.Upsert(ctx, config)
}

// GetServiceConfig gets a service configuration for a tenant
func (s *ServiceRegistry) GetServiceConfig(ctx context.Context, tenantID, serviceName string) (*domain.ServiceConfig, error) {
	config, err := s.repo.FindByTenantAndService(ctx, tenantID, serviceName)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, domain.ErrServiceNotFound
	}

	return config, nil
}

// GetTenantServices gets all service configurations for a tenant
func (s *ServiceRegistry) GetTenantServices(ctx context.Context, tenantID string) ([]*domain.ServiceConfig, error) {
	return s.repo.FindByTenant(ctx, tenantID)
}

// DeleteServiceConfig deletes a service configuration
func (s *ServiceRegistry) DeleteServiceConfig(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// CreateDefaultConfig creates or updates a default service configuration
func (s *ServiceRegistry) CreateDefaultConfig(ctx context.Context, config *domain.DefaultServiceConfig) error {
	return s.repo.UpsertDefaultConfig(ctx, config)
}

// GetDefaultConfig gets the default configuration for a service
func (s *ServiceRegistry) GetDefaultConfig(ctx context.Context, serviceName string) (*domain.DefaultServiceConfig, error) {
	config, err := s.repo.GetDefaultConfig(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, domain.ErrServiceNotFound
	}

	return config, nil
}

// GetAllDefaultConfigs gets all default service configurations
func (s *ServiceRegistry) GetAllDefaultConfigs(ctx context.Context) ([]*domain.DefaultServiceConfig, error) {
	return s.repo.GetAllDefaultConfigs(ctx)
}
