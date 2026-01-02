package service

import (
	"context"

	"github.com/vhvplatform/go-shared/errors"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"github.com/vhvplatform/go-tenant-service/internal/repository"
	"go.uber.org/zap"
)

// TenantService handles tenant business logic
type TenantService struct {
	tenantRepo     *repository.TenantRepository
	tenantUserRepo *repository.TenantUserRepository
	logger         *logger.Logger
}

// NewTenantService creates a new tenant service
func NewTenantService(
	tenantRepo *repository.TenantRepository,
	tenantUserRepo *repository.TenantUserRepository,
	log *logger.Logger,
) *TenantService {
	return &TenantService{
		tenantRepo:     tenantRepo,
		tenantUserRepo: tenantUserRepo,
		logger:         log,
	}
}

// CreateTenant creates a new tenant
func (s *TenantService) CreateTenant(ctx context.Context, req *domain.CreateTenantRequest) (*domain.Tenant, error) {
	// Check if tenant already exists
	existingTenant, err := s.tenantRepo.FindByName(ctx, req.Name)
	if err != nil {
		s.logger.Error("Failed to check existing tenant", zap.Error(err))
		return nil, errors.Internal("Failed to create tenant")
	}
	if existingTenant != nil {
		return nil, errors.Conflict("Tenant already exists with this name")
	}

	// Check domain if provided
	if req.Domain != "" {
		existingDomain, err := s.tenantRepo.FindByDomain(ctx, req.Domain)
		if err != nil {
			s.logger.Error("Failed to check existing domain", zap.Error(err))
			return nil, errors.Internal("Failed to create tenant")
		}
		if existingDomain != nil {
			return nil, errors.Conflict("Tenant already exists with this domain")
		}
	}

	// Create tenant
	tenant := &domain.Tenant{
		Name:             req.Name,
		Domain:           req.Domain,
		SubscriptionTier: req.SubscriptionTier,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		s.logger.Error("Failed to create tenant", zap.Error(err))
		return nil, errors.Internal("Failed to create tenant")
	}

	s.logger.Info("Tenant created successfully",
		zap.String("tenant_id", tenant.ID.Hex()),
		zap.String("name", tenant.Name),
	)

	return tenant, nil
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(ctx context.Context, id string) (*domain.Tenant, error) {
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get tenant", zap.String("tenant_id", id), zap.Error(err))
		return nil, errors.Internal("Failed to get tenant")
	}
	if tenant == nil {
		return nil, errors.NotFound("Tenant not found")
	}
	return tenant, nil
}

// ListTenants lists all tenants with pagination
func (s *TenantService) ListTenants(ctx context.Context, page, pageSize int) ([]*domain.Tenant, int64, error) {
	tenants, total, err := s.tenantRepo.List(ctx, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list tenants", zap.Error(err))
		return nil, 0, errors.Internal("Failed to list tenants")
	}
	return tenants, total, nil
}

// UpdateTenant updates a tenant
func (s *TenantService) UpdateTenant(ctx context.Context, id string, req *domain.UpdateTenantRequest) (*domain.Tenant, error) {
	// Get existing tenant
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to find tenant", zap.Error(err))
		return nil, errors.Internal("Failed to update tenant")
	}
	if tenant == nil {
		return nil, errors.NotFound("Tenant not found")
	}

	// Update fields
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Domain != "" {
		tenant.Domain = req.Domain
	}
	if req.SubscriptionTier != "" {
		tenant.SubscriptionTier = req.SubscriptionTier
	}

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		s.logger.Error("Failed to update tenant", zap.Error(err))
		return nil, errors.Internal("Failed to update tenant")
	}

	s.logger.Info("Tenant updated successfully",
		zap.String("tenant_id", tenant.ID.Hex()),
	)

	return tenant, nil
}

// DeleteTenant deletes a tenant
func (s *TenantService) DeleteTenant(ctx context.Context, id string) error {
	// Check if tenant exists
	tenant, err := s.tenantRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to find tenant", zap.Error(err))
		return errors.Internal("Failed to delete tenant")
	}
	if tenant == nil {
		return errors.NotFound("Tenant not found")
	}

	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete tenant", zap.Error(err))
		return errors.Internal("Failed to delete tenant")
	}

	s.logger.Info("Tenant deleted successfully",
		zap.String("tenant_id", id),
	)

	return nil
}

// AddUserToTenant adds a user to a tenant
func (s *TenantService) AddUserToTenant(ctx context.Context, tenantID, userID, role string) error {
	// Check if tenant exists
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to find tenant", zap.Error(err))
		return errors.Internal("Failed to add user to tenant")
	}
	if tenant == nil {
		return errors.NotFound("Tenant not found")
	}

	// Check if relationship already exists
	existing, err := s.tenantUserRepo.FindByTenantAndUser(ctx, tenantID, userID)
	if err != nil {
		s.logger.Error("Failed to check existing relationship", zap.Error(err))
		return errors.Internal("Failed to add user to tenant")
	}
	if existing != nil {
		return errors.Conflict("User already belongs to this tenant")
	}

	// Add user to tenant
	tenantUser := &domain.TenantUser{
		TenantID: tenantID,
		UserID:   userID,
		Role:     role,
	}

	if err := s.tenantUserRepo.AddUser(ctx, tenantUser); err != nil {
		s.logger.Error("Failed to add user to tenant", zap.Error(err))
		return errors.Internal("Failed to add user to tenant")
	}

	s.logger.Info("User added to tenant successfully",
		zap.String("tenant_id", tenantID),
		zap.String("user_id", userID),
	)

	return nil
}

// RemoveUserFromTenant removes a user from a tenant
func (s *TenantService) RemoveUserFromTenant(ctx context.Context, tenantID, userID string) error {
	// Check if relationship exists
	existing, err := s.tenantUserRepo.FindByTenantAndUser(ctx, tenantID, userID)
	if err != nil {
		s.logger.Error("Failed to check existing relationship", zap.Error(err))
		return errors.Internal("Failed to remove user from tenant")
	}
	if existing == nil {
		return errors.NotFound("User not found in tenant")
	}

	if err := s.tenantUserRepo.RemoveUser(ctx, tenantID, userID); err != nil {
		s.logger.Error("Failed to remove user from tenant", zap.Error(err))
		return errors.Internal("Failed to remove user from tenant")
	}

	s.logger.Info("User removed from tenant successfully",
		zap.String("tenant_id", tenantID),
		zap.String("user_id", userID),
	)

	return nil
}

// GetTenantUsers retrieves all users in a tenant
func (s *TenantService) GetTenantUsers(ctx context.Context, tenantID string) ([]*domain.TenantUser, error) {
	// Check if tenant exists
	tenant, err := s.tenantRepo.FindByID(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to find tenant", zap.Error(err))
		return nil, errors.Internal("Failed to get tenant users")
	}
	if tenant == nil {
		return nil, errors.NotFound("Tenant not found")
	}

	tenantUsers, err := s.tenantUserRepo.ListUsersByTenant(ctx, tenantID)
	if err != nil {
		s.logger.Error("Failed to get tenant users", zap.Error(err))
		return nil, errors.Internal("Failed to get tenant users")
	}

	return tenantUsers, nil
}
