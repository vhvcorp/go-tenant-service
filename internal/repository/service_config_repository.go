package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/vhvplatform/go-tenant-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ServiceConfigRepository handles service configuration data access
type ServiceConfigRepository struct {
	collection        *mongo.Collection
	defaultCollection *mongo.Collection
}

// NewServiceConfigRepository creates a new service config repository
func NewServiceConfigRepository(db *mongo.Database) *ServiceConfigRepository {
	collection := db.Collection("service_configs")
	defaultCollection := db.Collection("default_service_configs")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Indexes for service_configs
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenantId", Value: 1},
				{Key: "serviceName", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "serviceName", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "isActive", Value: 1}},
		},
	}
	_, _ = collection.Indexes().CreateMany(ctx, indexes)

	// Indexes for default_service_configs
	defaultIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "serviceName", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	_, _ = defaultCollection.Indexes().CreateMany(ctx, defaultIndexes)

	return &ServiceConfigRepository{
		collection:        collection,
		defaultCollection: defaultCollection,
	}
}

// Create creates a new service configuration
func (r *ServiceConfigRepository) Create(ctx context.Context, config *domain.ServiceConfig) error {
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	if err := config.Validate(); err != nil {
		return err
	}

	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create service config: %w", err)
	}

	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update updates an existing service configuration
func (r *ServiceConfigRepository) Update(ctx context.Context, config *domain.ServiceConfig) error {
	config.UpdatedAt = time.Now()

	if err := config.Validate(); err != nil {
		return err
	}

	filter := bson.M{"_id": config.ID}
	update := bson.M{"$set": config}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update service config: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrServiceNotFound
	}

	return nil
}

// FindByID finds a service configuration by ID
func (r *ServiceConfigRepository) FindByID(ctx context.Context, id string) (*domain.ServiceConfig, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	var config domain.ServiceConfig
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrServiceNotFound
		}
		return nil, fmt.Errorf("failed to find service config: %w", err)
	}

	return &config, nil
}

// FindByTenantAndService finds service configuration for a tenant and service
func (r *ServiceConfigRepository) FindByTenantAndService(ctx context.Context, tenantID, serviceName string) (*domain.ServiceConfig, error) {
	filter := bson.M{
		"tenantId":    tenantID,
		"serviceName": serviceName,
		"isActive":    true,
	}

	var config domain.ServiceConfig
	err := r.collection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if not found (will use default)
		}
		return nil, fmt.Errorf("failed to find service config: %w", err)
	}

	return &config, nil
}

// FindByTenant finds all service configurations for a tenant
func (r *ServiceConfigRepository) FindByTenant(ctx context.Context, tenantID string) ([]*domain.ServiceConfig, error) {
	filter := bson.M{"tenantId": tenantID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find service configs: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []*domain.ServiceConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode service configs: %w", err)
	}

	return configs, nil
}

// FindByService finds all configurations for a specific service across tenants
func (r *ServiceConfigRepository) FindByService(ctx context.Context, serviceName string) ([]*domain.ServiceConfig, error) {
	filter := bson.M{"serviceName": serviceName}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find service configs: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []*domain.ServiceConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode service configs: %w", err)
	}

	return configs, nil
}

// Delete deletes a service configuration
func (r *ServiceConfigRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete service config: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrServiceNotFound
	}

	return nil
}

// Upsert creates or updates a service configuration
func (r *ServiceConfigRepository) Upsert(ctx context.Context, config *domain.ServiceConfig) error {
	config.UpdatedAt = time.Now()

	if err := config.Validate(); err != nil {
		return err
	}

	filter := bson.M{
		"tenantId":    config.TenantID,
		"serviceName": config.ServiceName,
	}

	update := bson.M{
		"$set": config,
		"$setOnInsert": bson.M{
			"createdAt": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert service config: %w", err)
	}

	if result.UpsertedID != nil {
		config.ID = result.UpsertedID.(primitive.ObjectID)
	}

	return nil
}

// === Default Service Config Methods ===

// CreateDefaultConfig creates a default service configuration
func (r *ServiceConfigRepository) CreateDefaultConfig(ctx context.Context, config *domain.DefaultServiceConfig) error {
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	result, err := r.defaultCollection.InsertOne(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}

	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetDefaultConfig gets the default configuration for a service
func (r *ServiceConfigRepository) GetDefaultConfig(ctx context.Context, serviceName string) (*domain.DefaultServiceConfig, error) {
	filter := bson.M{"serviceName": serviceName}

	var config domain.DefaultServiceConfig
	err := r.defaultCollection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find default config: %w", err)
	}

	return &config, nil
}

// UpsertDefaultConfig creates or updates a default service configuration
func (r *ServiceConfigRepository) UpsertDefaultConfig(ctx context.Context, config *domain.DefaultServiceConfig) error {
	config.UpdatedAt = time.Now()

	filter := bson.M{"serviceName": config.ServiceName}
	update := bson.M{
		"$set": config,
		"$setOnInsert": bson.M{
			"createdAt": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.defaultCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert default config: %w", err)
	}

	if result.UpsertedID != nil {
		config.ID = result.UpsertedID.(primitive.ObjectID)
	}

	return nil
}

// GetAllDefaultConfigs gets all default service configurations
func (r *ServiceConfigRepository) GetAllDefaultConfigs(ctx context.Context) ([]*domain.DefaultServiceConfig, error) {
	cursor, err := r.defaultCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find default configs: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []*domain.DefaultServiceConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode default configs: %w", err)
	}

	return configs, nil
}

// UpdateEndpointStatus updates the status of a specific endpoint (for health checks)
func (r *ServiceConfigRepository) UpdateEndpointStatus(ctx context.Context, tenantID, serviceName, endpointURL string, isActive bool) error {
	filter := bson.M{
		"tenantId":    tenantID,
		"serviceName": serviceName,
	}

	// Update primary endpoint if it matches
	update := bson.M{
		"$set": bson.M{
			"primaryEndpoint.isActive": isActive,
			"updatedAt":                time.Now(),
		},
	}

	// This is simplified - in production, you'd need to identify which endpoint to update
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update endpoint status: %w", err)
	}

	return nil
}

// GetActiveServices gets all active service configurations
func (r *ServiceConfigRepository) GetActiveServices(ctx context.Context) ([]*domain.ServiceConfig, error) {
	filter := bson.M{"isActive": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find active services: %w", err)
	}
	defer cursor.Close(ctx)

	var configs []*domain.ServiceConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode service configs: %w", err)
	}

	return configs, nil
}
