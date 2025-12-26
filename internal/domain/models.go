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
	SubscriptionTier string                 `bson:"subscription_tier" json:"subscription_tier"`
	IsActive         bool                   `bson:"is_active" json:"is_active"`
	Settings         map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
	CreatedAt        time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time              `bson:"updated_at" json:"updated_at"`
}

// TenantUser represents a user's relationship with a tenant
type TenantUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  string             `bson:"tenant_id" json:"tenant_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Role      string             `bson:"role" json:"role"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Subscription tiers
const (
	SubscriptionFree         = "free"
	SubscriptionBasic        = "basic"
	SubscriptionProfessional = "professional"
	SubscriptionEnterprise   = "enterprise"
)
