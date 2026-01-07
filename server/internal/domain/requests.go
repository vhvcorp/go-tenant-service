package domain

// CreateTenantRequest represents a tenant creation request
type CreateTenantRequest struct {
	Name             string `json:"name" binding:"required"`
	Domain           string `json:"domain"`
	SubscriptionTier string `json:"subscription_tier"`
}

// UpdateTenantRequest represents a tenant update request
type UpdateTenantRequest struct {
	Name             string `json:"name"`
	Domain           string `json:"domain"`
	SubscriptionTier string `json:"subscription_tier"`
}

// AddUserRequest represents adding a user to tenant
type AddUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

// ListTenantsRequest represents a list tenants request
type ListTenantsRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

// TenantResponse represents a tenant in API responses
type TenantResponse struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Domain           string                 `json:"domain,omitempty"`
	SubscriptionTier string                 `json:"subscription_tier"`
	IsActive         bool                   `json:"is_active"`
	Settings         map[string]interface{} `json:"settings,omitempty"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

// ListTenantsResponse represents a paginated list of tenants
type ListTenantsResponse struct {
	Tenants  []TenantResponse `json:"tenants"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}
