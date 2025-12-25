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

// TenantRepository handles tenant data access
type TenantRepository struct {
	collection *mongo.Collection
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *mongo.Database) *TenantRepository {
	collection := db.Collection("tenants")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "domain", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}
	
	_, _ = collection.Indexes().CreateMany(ctx, indexes)
	
	return &TenantRepository{collection: collection}
}

// Create creates a new tenant
func (r *TenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	tenant.IsActive = true
	
	if tenant.SubscriptionTier == "" {
		tenant.SubscriptionTier = domain.SubscriptionFree
	}
	
	result, err := r.collection.InsertOne(ctx, tenant)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	
	tenant.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a tenant by ID
func (r *TenantRepository) FindByID(ctx context.Context, id string) (*domain.Tenant, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}
	
	var tenant domain.Tenant
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tenant: %w", err)
	}
	return &tenant, nil
}

// FindByName finds a tenant by name
func (r *TenantRepository) FindByName(ctx context.Context, name string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tenant by name: %w", err)
	}
	return &tenant, nil
}

// FindByDomain finds a tenant by domain
func (r *TenantRepository) FindByDomain(ctx context.Context, domainName string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	err := r.collection.FindOne(ctx, bson.M{"domain": domainName}).Decode(&tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tenant by domain: %w", err)
	}
	return &tenant, nil
}

// List lists all tenants with pagination
func (r *TenantRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Tenant, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	skip := (page - 1) * pageSize
	
	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tenants: %w", err)
	}
	
	// Get tenants
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer cursor.Close(ctx)
	
	var tenants []*domain.Tenant
	if err := cursor.All(ctx, &tenants); err != nil {
		return nil, 0, fmt.Errorf("failed to decode tenants: %w", err)
	}
	
	return tenants, total, nil
}

// Update updates a tenant
func (r *TenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	tenant.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": tenant.ID},
		bson.M{"$set": tenant},
	)
	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}

// Delete soft deletes a tenant
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid tenant ID: %w", err)
	}
	
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}
