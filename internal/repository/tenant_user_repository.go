package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vhvcorp/go-tenant-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TenantUserRepository handles tenant-user relationship data access
type TenantUserRepository struct {
	collection *mongo.Collection
}

// NewTenantUserRepository creates a new tenant-user repository
func NewTenantUserRepository(db *mongo.Database) *TenantUserRepository {
	collection := db.Collection("tenant_users")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "tenant_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	}
	
	_, _ = collection.Indexes().CreateMany(ctx, indexes)
	
	return &TenantUserRepository{collection: collection}
}

// AddUser adds a user to a tenant
func (r *TenantUserRepository) AddUser(ctx context.Context, tenantUser *domain.TenantUser) error {
	tenantUser.CreatedAt = time.Now()
	tenantUser.UpdatedAt = time.Now()
	tenantUser.IsActive = true
	
	result, err := r.collection.InsertOne(ctx, tenantUser)
	if err != nil {
		return fmt.Errorf("failed to add user to tenant: %w", err)
	}
	
	tenantUser.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// RemoveUser removes a user from a tenant
func (r *TenantUserRepository) RemoveUser(ctx context.Context, tenantID, userID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"tenant_id": tenantID,
			"user_id":   userID,
		},
		bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to remove user from tenant: %w", err)
	}
	return nil
}

// FindByTenantAndUser finds a tenant-user relationship
func (r *TenantUserRepository) FindByTenantAndUser(ctx context.Context, tenantID, userID string) (*domain.TenantUser, error) {
	var tenantUser domain.TenantUser
	err := r.collection.FindOne(ctx, bson.M{
		"tenant_id": tenantID,
		"user_id":   userID,
		"is_active": true,
	}).Decode(&tenantUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tenant-user: %w", err)
	}
	return &tenantUser, nil
}

// ListUsersByTenant lists all users in a tenant
func (r *TenantUserRepository) ListUsersByTenant(ctx context.Context, tenantID string) ([]*domain.TenantUser, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"tenant_id": tenantID,
		"is_active": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tenant users: %w", err)
	}
	defer cursor.Close(ctx)
	
	var tenantUsers []*domain.TenantUser
	if err := cursor.All(ctx, &tenantUsers); err != nil {
		return nil, fmt.Errorf("failed to decode tenant users: %w", err)
	}
	
	return tenantUsers, nil
}

// ListTenantsByUser lists all tenants for a user
func (r *TenantUserRepository) ListTenantsByUser(ctx context.Context, userID string) ([]*domain.TenantUser, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":   userID,
		"is_active": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list user tenants: %w", err)
	}
	defer cursor.Close(ctx)
	
	var tenantUsers []*domain.TenantUser
	if err := cursor.All(ctx, &tenantUsers); err != nil {
		return nil, fmt.Errorf("failed to decode user tenants: %w", err)
	}
	
	return tenantUsers, nil
}
