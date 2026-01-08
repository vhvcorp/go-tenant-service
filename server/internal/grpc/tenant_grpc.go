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
	tenantService *service.TenantService
	logger        *logger.Logger
}

// NewTenantServiceServer creates a new gRPC tenant service server
func NewTenantServiceServer(tenantService *service.TenantService, log *logger.Logger) *TenantServiceServer {
	return &TenantServiceServer{
		tenantService: tenantService,
		logger:        log,
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
