package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tenant represents an organization/tenant in the system
type Tenant struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name             string                 `bson:"name" json:"name"`
	Domain           string                 `bson:"domain,omitempty" json:"domain,omitempty"`
	SubscriptionTier string                 `bson:"subscriptionTier" json:"subscription_tier"`
	IsActive         bool                   `bson:"isActive" json:"is_active"`
	AuthSettings     AuthSettings           `bson:"authSettings" json:"auth_settings"`
	DefaultService   string                 `bson:"defaultService" json:"default_service"`
	Settings         map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
	CreatedAt        time.Time              `bson:"createdAt" json:"created_at"`
	UpdatedAt        time.Time              `bson:"updatedAt" json:"updated_at"`
}

// AuthSettings defines authentication configuration for a tenant
type AuthSettings struct {
	AllowedLoginMethods []string `bson:"allowedLoginMethods" json:"allowed_login_methods"`
}

// User represents a global user account
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username,omitempty" json:"username,omitempty"`
	Email          string             `bson:"email,omitempty" json:"email,omitempty"`
	Phone          string             `bson:"phone,omitempty" json:"phone,omitempty"`
	DocumentNumber string             `bson:"documentNumber,omitempty" json:"document_number,omitempty"`
	PasswordHash   string             `bson:"passwordHash" json:"-"`
	IsActive       bool               `bson:"isActive" json:"is_active"`
	CreatedAt      time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updated_at"`
}

// TenantUser represents a user's relationship with a tenant
type TenantUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  string             `bson:"tenantId" json:"tenant_id"`
	UserID    string             `bson:"userId" json:"user_id"`
	Role      string             `bson:"role" json:"role"`
	IsActive  bool               `bson:"isActive" json:"is_active"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"`
}

// Subscription tiers
const (
	SubscriptionFree         = "free"
	SubscriptionBasic        = "basic"
	SubscriptionProfessional = "professional"
	SubscriptionEnterprise   = "enterprise"
)
