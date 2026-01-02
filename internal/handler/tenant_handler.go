package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/errors"
	"github.com/vhvplatform/go-shared/logger"
	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"github.com/vhvplatform/go-tenant-service/internal/service"
	"go.uber.org/zap"
)

// TenantHandler handles HTTP requests for tenants
type TenantHandler struct {
	tenantService *service.TenantService
	logger        *logger.Logger
}

// TenantUserResponse represents a tenant user in API responses
type TenantUserResponse struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService *service.TenantService, log *logger.Logger) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
		logger:        log,
	}
}

// CreateTenant godoc
// @Summary Create a new tenant
// @Description Create a new tenant in the system
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body domain.CreateTenantRequest true "Tenant creation request"
// @Success 201 {object} map[string]interface{} "Tenant created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 409 {object} map[string]interface{} "Tenant already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req domain.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	tenant, err := h.tenantService.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": h.toTenantResponse(tenant)})
}

// GetTenant godoc
// @Summary Get tenant by ID
// @Description Get tenant details by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "Tenant details"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants/{id} [get]
func (h *TenantHandler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")

	tenant, err := h.tenantService.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.toTenantResponse(tenant)})
}

// ListTenants godoc
// @Summary List all tenants
// @Description Get paginated list of tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{} "List of tenants"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants [get]
func (h *TenantHandler) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tenants, total, err := h.tenantService.ListTenants(c.Request.Context(), page, pageSize)
	if err != nil {
		h.respondError(c, err)
		return
	}

	tenantResponses := make([]domain.TenantResponse, len(tenants))
	for i, tenant := range tenants {
		tenantResponses[i] = h.toTenantResponse(tenant)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": domain.ListTenantsResponse{
			Tenants:  tenantResponses,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// UpdateTenant godoc
// @Summary Update tenant
// @Description Update tenant information
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body domain.UpdateTenantRequest true "Tenant update request"
// @Success 200 {object} map[string]interface{} "Tenant updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	tenantID := c.Param("id")

	var req domain.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	tenant, err := h.tenantService.UpdateTenant(c.Request.Context(), tenantID, &req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": h.toTenantResponse(tenant)})
}

// DeleteTenant godoc
// @Summary Delete tenant (soft delete)
// @Description Soft delete a tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "Tenant deleted"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	tenantID := c.Param("id")

	if err := h.tenantService.DeleteTenant(c.Request.Context(), tenantID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tenant deleted successfully"})
}

// AddUserToTenant godoc
// @Summary Add user to tenant
// @Description Add a user to a tenant with specific role
// @Tags tenant-users
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param user body domain.AddUserRequest true "User addition request"
// @Success 200 {object} map[string]interface{} "User added to tenant"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants/{id}/users [post]
func (h *TenantHandler) AddUserToTenant(c *gin.Context) {
	tenantID := c.Param("id")

	var req domain.AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondError(c, errors.BadRequest("Invalid request body"))
		return
	}

	if err := h.tenantService.AddUserToTenant(c.Request.Context(), tenantID, req.UserID, req.Role); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to tenant successfully"})
}

// RemoveUserFromTenant godoc
// @Summary Remove user from tenant
// @Description Remove a user from a tenant
// @Tags tenant-users
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User removed from tenant"
// @Failure 404 {object} map[string]interface{} "User or tenant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants/{id}/users/{user_id} [delete]
func (h *TenantHandler) RemoveUserFromTenant(c *gin.Context) {
	tenantID := c.Param("id")
	userID := c.Param("user_id")

	if err := h.tenantService.RemoveUserFromTenant(c.Request.Context(), tenantID, userID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from tenant successfully"})
}

// GetTenantUsers godoc
// @Summary Get tenant users
// @Description Get list of users in a tenant
// @Tags tenant-users
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "List of tenant users"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/tenants/{id}/users [get]
func (h *TenantHandler) GetTenantUsers(c *gin.Context) {
	tenantID := c.Param("id")

	tenantUsers, err := h.tenantService.GetTenantUsers(c.Request.Context(), tenantID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	userResponses := make([]TenantUserResponse, len(tenantUsers))
	for i, tu := range tenantUsers {
		userResponses[i] = TenantUserResponse{
			ID:        tu.ID.Hex(),
			TenantID:  tu.TenantID,
			UserID:    tu.UserID,
			Role:      tu.Role,
			IsActive:  tu.IsActive,
			CreatedAt: tu.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: tu.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": userResponses})
}

// toTenantResponse converts a tenant domain model to a response
func (h *TenantHandler) toTenantResponse(tenant *domain.Tenant) domain.TenantResponse {
	return domain.TenantResponse{
		ID:               tenant.ID.Hex(),
		Name:             tenant.Name,
		Domain:           tenant.Domain,
		SubscriptionTier: tenant.SubscriptionTier,
		IsActive:         tenant.IsActive,
		Settings:         tenant.Settings,
		CreatedAt:        tenant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        tenant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// respondError responds with an error
func (h *TenantHandler) respondError(c *gin.Context, err error) {
	appErr := errors.FromError(err)
	h.logger.Error("Request failed",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("error", appErr.Message),
	)
	c.JSON(appErr.StatusCode, gin.H{"error": appErr})
}
