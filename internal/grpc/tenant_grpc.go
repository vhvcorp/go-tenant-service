package grpc

import (
	"context"
	"time"

	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"github.com/vhvplatform/go-tenant-service/internal/service"
	pb "github.com/vhvplatform/go-tenant-service/proto"
	"go.uber.org/zap"
)

// TenantServiceServer implements the gRPC tenant service
type TenantServiceServer struct {
	pb.UnimplementedTenantServiceServer
	tenantService   *service.TenantService
	registryService *service.ServiceRegistry
	logger          *logger.Logger
}

// NewTenantServiceServer creates a new gRPC tenant service server
func NewTenantServiceServer(tenantService *service.TenantService, registryService *service.ServiceRegistry, log *logger.Logger) *TenantServiceServer {
	return &TenantServiceServer{
		tenantService:   tenantService,
		registryService: registryService,
		logger:          log,
	}
}

// GetTenant retrieves a tenant by ID
func (s *TenantServiceServer) GetTenant(ctx context.Context, req *pb.GetTenantRequest) (*pb.GetTenantResponse, error) {
	tenant, err := s.tenantService.GetTenant(ctx, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to get tenant", zap.Error(err))
		return nil, err
	}

	return &pb.GetTenantResponse{
		Tenant: s.toProtoTenant(tenant),
	}, nil
}

// ListTenants lists all tenants
func (s *TenantServiceServer) ListTenants(ctx context.Context, req *pb.ListTenantsRequest) (*pb.ListTenantsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	tenants, total, err := s.tenantService.ListTenants(ctx, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list tenants", zap.Error(err))
		return nil, err
	}

	protoTenants := make([]*pb.Tenant, len(tenants))
	for i, tenant := range tenants {
		protoTenants[i] = s.toProtoTenant(tenant)
	}

	return &pb.ListTenantsResponse{
		Tenants: protoTenants,
		Total:   int32(total),
	}, nil
}

// CreateTenant creates a new tenant
func (s *TenantServiceServer) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.CreateTenantResponse, error) {
	createReq := &domain.CreateTenantRequest{
		Name:             req.Name,
		Domain:           req.Domain,
		SubscriptionTier: req.SubscriptionTier,
	}

	tenant, err := s.tenantService.CreateTenant(ctx, createReq)
	if err != nil {
		s.logger.Error("Failed to create tenant", zap.Error(err))
		return nil, err
	}

	return &pb.CreateTenantResponse{
		Tenant: s.toProtoTenant(tenant),
	}, nil
}

// UpdateTenant updates a tenant
func (s *TenantServiceServer) UpdateTenant(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.UpdateTenantResponse, error) {
	updateReq := &domain.UpdateTenantRequest{
		Name:             req.Name,
		Domain:           req.Domain,
		SubscriptionTier: req.SubscriptionTier,
	}

	tenant, err := s.tenantService.UpdateTenant(ctx, req.TenantId, updateReq)
	if err != nil {
		s.logger.Error("Failed to update tenant", zap.Error(err))
		return nil, err
	}

	return &pb.UpdateTenantResponse{
		Tenant: s.toProtoTenant(tenant),
	}, nil
}

// DeleteTenant deletes a tenant
func (s *TenantServiceServer) DeleteTenant(ctx context.Context, req *pb.DeleteTenantRequest) (*pb.DeleteTenantResponse, error) {
	err := s.tenantService.DeleteTenant(ctx, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to delete tenant", zap.Error(err))
		return nil, err
	}

	return &pb.DeleteTenantResponse{
		Success: true,
	}, nil
}

// AddUserToTenant adds a user to a tenant
func (s *TenantServiceServer) AddUserToTenant(ctx context.Context, req *pb.AddUserToTenantRequest) (*pb.AddUserToTenantResponse, error) {
	err := s.tenantService.AddUserToTenant(ctx, req.TenantId, req.UserId, req.Role)
	if err != nil {
		s.logger.Error("Failed to add user to tenant", zap.Error(err))
		return nil, err
	}

	return &pb.AddUserToTenantResponse{
		Success: true,
	}, nil
}

// RemoveUserFromTenant removes a user from a tenant
func (s *TenantServiceServer) RemoveUserFromTenant(ctx context.Context, req *pb.RemoveUserFromTenantRequest) (*pb.RemoveUserFromTenantResponse, error) {
	err := s.tenantService.RemoveUserFromTenant(ctx, req.TenantId, req.UserId)
	if err != nil {
		s.logger.Error("Failed to remove user from tenant", zap.Error(err))
		return nil, err
	}

	return &pb.RemoveUserFromTenantResponse{
		Success: true,
	}, nil
}

func (s *TenantServiceServer) toProtoTenant(tenant *domain.Tenant) *pb.Tenant {
	return &pb.Tenant{
		Id:               tenant.ID.Hex(),
		Name:             tenant.Name,
		Domain:           tenant.Domain,
		SubscriptionTier: tenant.SubscriptionTier,
		IsActive:         tenant.IsActive,
		CreatedAt:        tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        tenant.UpdatedAt.Format(time.RFC3339),
	}
}

// === Service Registry Handlers ===

// GetServiceConfig gets service configuration for a tenant
func (s *TenantServiceServer) GetServiceConfig(ctx context.Context, req *pb.GetServiceConfigRequest) (*pb.GetServiceConfigResponse, error) {
	config, err := s.registryService.GetServiceConfig(ctx, req.TenantId, req.ServiceName)
	if err != nil {
		s.logger.Error("Failed to get service config", zap.Error(err))
		return nil, err
	}

	return &pb.GetServiceConfigResponse{
		Config: s.toProtoServiceConfig(config),
	}, nil
}

// UpdateServiceConfig creates or updates service configuration for a tenant
func (s *TenantServiceServer) UpdateServiceConfig(ctx context.Context, req *pb.UpdateServiceConfigRequest) (*pb.UpdateServiceConfigResponse, error) {
	config := s.fromProtoServiceConfig(req.Config)
	config.TenantID = req.TenantId
	config.ServiceName = req.ServiceName

	err := s.registryService.CreateOrUpdateServiceConfig(ctx, config)
	if err != nil {
		s.logger.Error("Failed to update service config", zap.Error(err))
		return nil, err
	}

	return &pb.UpdateServiceConfigResponse{
		Config: s.toProtoServiceConfig(config),
	}, nil
}

// GetServiceURL resolves the service URL for a tenant
func (s *TenantServiceServer) GetServiceURL(ctx context.Context, req *pb.GetServiceURLRequest) (*pb.GetServiceURLResponse, error) {
	result, err := s.registryService.GetServiceURL(ctx, req.TenantId, req.ServiceName)
	if err != nil {
		s.logger.Error("Failed to get service URL", zap.Error(err))
		return &pb.GetServiceURLResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.GetServiceURLResponse{
		Url:           result.ResolvedURL,
		IsDefault:     result.IsDefault,
		Success:       result.Success,
		Error:         result.Error,
		AttemptedUrls: result.AttemptedURLs,
	}, nil
}

// ListTenantServices lists all service configurations for a tenant
func (s *TenantServiceServer) ListTenantServices(ctx context.Context, req *pb.ListTenantServicesRequest) (*pb.ListTenantServicesResponse, error) {
	configs, err := s.registryService.GetTenantServices(ctx, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to list tenant services", zap.Error(err))
		return nil, err
	}

	protoConfigs := make([]*pb.ServiceConfig, len(configs))
	for i, config := range configs {
		protoConfigs[i] = s.toProtoServiceConfig(config)
	}

	return &pb.ListTenantServicesResponse{
		Services: protoConfigs,
	}, nil
}

// GetServiceHealth gets health status for a service
func (s *TenantServiceServer) GetServiceHealth(ctx context.Context, req *pb.GetServiceHealthRequest) (*pb.GetServiceHealthResponse, error) {
	config, err := s.registryService.GetServiceConfig(ctx, req.TenantId, req.ServiceName)
	if err != nil {
		s.logger.Error("Failed to get service config for health", zap.Error(err))
		return nil, err
	}

	healths := []*pb.ServiceHealth{}

	// Primary endpoint health
	primaryStatus := s.registryService.GetHealthStatus(req.TenantId, req.ServiceName, config.PrimaryEndpoint.URL)
	healths = append(healths, s.toProtoServiceHealth(primaryStatus))

	// Fallback chain health
	for _, endpoint := range config.FallbackChain {
		status := s.registryService.GetHealthStatus(req.TenantId, req.ServiceName, endpoint.URL)
		healths = append(healths, s.toProtoServiceHealth(status))
	}

	return &pb.GetServiceHealthResponse{
		Healths: healths,
	}, nil
}

// === Proto Conversion Helpers ===

func (s *TenantServiceServer) toProtoServiceConfig(config *domain.ServiceConfig) *pb.ServiceConfig {
	fallbackChain := make([]*pb.ServiceEndpoint, len(config.FallbackChain))
	for i, endpoint := range config.FallbackChain {
		fallbackChain[i] = s.toProtoServiceEndpoint(&endpoint)
	}

	return &pb.ServiceConfig{
		Id:                  config.ID.Hex(),
		TenantId:            config.TenantID,
		ServiceName:         config.ServiceName,
		PrimaryEndpoint:     s.toProtoServiceEndpoint(&config.PrimaryEndpoint),
		FallbackChain:       fallbackChain,
		DefaultServiceUrl:   config.DefaultServiceURL,
		HealthCheck:         s.toProtoHealthCheck(&config.HealthCheck),
		LoadBalanceStrategy: config.LoadBalanceStrategy,
		IsActive:            config.IsActive,
		Metadata:            config.Metadata,
		CreatedAt:           config.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           config.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *TenantServiceServer) toProtoServiceEndpoint(endpoint *domain.ServiceEndpoint) *pb.ServiceEndpoint {
	return &pb.ServiceEndpoint{
		Url:      endpoint.URL,
		Priority: int32(endpoint.Priority),
		Weight:   int32(endpoint.Weight),
		Timeout:  int32(endpoint.Timeout),
		Headers:  endpoint.Headers,
		IsActive: endpoint.IsActive,
	}
}

func (s *TenantServiceServer) toProtoHealthCheck(config *domain.HealthCheckConfig) *pb.HealthCheckConfig {
	return &pb.HealthCheckConfig{
		Enabled:       config.Enabled,
		Path:          config.Path,
		Method:        config.Method,
		Interval:      int32(config.Interval),
		Timeout:       int32(config.Timeout),
		FailThreshold: int32(config.FailThreshold),
	}
}

func (s *TenantServiceServer) toProtoServiceHealth(status *domain.ServiceStatus) *pb.ServiceHealth {
	return &pb.ServiceHealth{
		EndpointUrl:      status.EndpointURL,
		IsHealthy:        status.IsHealthy,
		LastChecked:      status.LastChecked.Format(time.RFC3339),
		LastSuccessful:   status.LastSuccessful.Format(time.RFC3339),
		LastFailure:      status.LastFailure.Format(time.RFC3339),
		ConsecutiveFails: int32(status.ConsecutiveFails),
		LastError:        status.LastError,
	}
}

func (s *TenantServiceServer) fromProtoServiceConfig(proto *pb.ServiceConfig) *domain.ServiceConfig {
	fallbackChain := make([]domain.ServiceEndpoint, len(proto.FallbackChain))
	for i, endpoint := range proto.FallbackChain {
		fallbackChain[i] = *s.fromProtoServiceEndpoint(endpoint)
	}

	return &domain.ServiceConfig{
		TenantID:            proto.TenantId,
		ServiceName:         proto.ServiceName,
		PrimaryEndpoint:     *s.fromProtoServiceEndpoint(proto.PrimaryEndpoint),
		FallbackChain:       fallbackChain,
		DefaultServiceURL:   proto.DefaultServiceUrl,
		HealthCheck:         *s.fromProtoHealthCheck(proto.HealthCheck),
		LoadBalanceStrategy: proto.LoadBalanceStrategy,
		IsActive:            proto.IsActive,
		Metadata:            proto.Metadata,
	}
}

func (s *TenantServiceServer) fromProtoServiceEndpoint(proto *pb.ServiceEndpoint) *domain.ServiceEndpoint {
	return &domain.ServiceEndpoint{
		URL:      proto.Url,
		Priority: int(proto.Priority),
		Weight:   int(proto.Weight),
		Timeout:  int(proto.Timeout),
		Headers:  proto.Headers,
		IsActive: proto.IsActive,
	}
}

func (s *TenantServiceServer) fromProtoHealthCheck(proto *pb.HealthCheckConfig) *domain.HealthCheckConfig {
	return &domain.HealthCheckConfig{
		Enabled:       proto.Enabled,
		Path:          proto.Path,
		Method:        proto.Method,
		Interval:      int(proto.Interval),
		Timeout:       int(proto.Timeout),
		FailThreshold: int(proto.FailThreshold),
	}
}
