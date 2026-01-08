// Migration: 002_service_config
// Description: Setup service_configs and default_service_configs collections with indexes
// Date: 2025-01-10

db = db.getSiblingDB('tenant_service');

// Create service_configs collection
db.createCollection('service_configs');

// Create indexes for service_configs
db.service_configs.createIndex(
    { tenantId: 1, serviceName: 1 },
    { unique: true, name: 'idx_tenant_service_unique' }
);

db.service_configs.createIndex(
    { serviceName: 1 },
    { name: 'idx_service_name' }
);

db.service_configs.createIndex(
    { isActive: 1 },
    { name: 'idx_is_active' }
);

db.service_configs.createIndex(
    { createdAt: -1 },
    { name: 'idx_created_at' }
);

// Create default_service_configs collection
db.createCollection('default_service_configs');

// Create indexes for default_service_configs
db.default_service_configs.createIndex(
    { serviceName: 1 },
    { unique: true, name: 'idx_service_name_unique' }
);

// Insert default service configurations
db.default_service_configs.insertMany([
    {
        serviceName: 'auth',
        defaultURL: 'http://go-auth-service:50051',
        description: 'Authentication and Authorization Service',
        healthCheck: {
            enabled: true,
            path: '/health',
            method: 'GET',
            interval: 30,
            timeout: 5,
            failThreshold: 3
        },
        createdAt: new Date(),
        updatedAt: new Date()
    },
    {
        serviceName: 'user',
        defaultURL: 'http://go-user-service:50052',
        description: 'User Management Service',
        healthCheck: {
            enabled: true,
            path: '/health',
            method: 'GET',
            interval: 30,
            timeout: 5,
            failThreshold: 3
        },
        createdAt: new Date(),
        updatedAt: new Date()
    },
    {
        serviceName: 'tenant',
        defaultURL: 'http://go-tenant-service:50053',
        description: 'Tenant Management Service',
        healthCheck: {
            enabled: true,
            path: '/health',
            method: 'GET',
            interval: 30,
            timeout: 5,
            failThreshold: 3
        },
        createdAt: new Date(),
        updatedAt: new Date()
    },
    {
        serviceName: 'notification',
        defaultURL: 'http://go-notification-service:50054',
        description: 'Notification Service',
        healthCheck: {
            enabled: true,
            path: '/health',
            method: 'GET',
            interval: 30,
            timeout: 5,
            failThreshold: 3
        },
        createdAt: new Date(),
        updatedAt: new Date()
    },
    {
        serviceName: 'cms',
        defaultURL: 'http://go-cms-service:50055',
        description: 'Content Management Service',
        healthCheck: {
            enabled: true,
            path: '/health',
            method: 'GET',
            interval: 30,
            timeout: 5,
            failThreshold: 3
        },
        createdAt: new Date(),
        updatedAt: new Date()
    },
    {
        serviceName: 'config',
        defaultURL: 'http://go-system-config-service:50056',
        description: 'System Configuration Service',
        healthCheck: {
            enabled: true,
            path: '/health',
            method: 'GET',
            interval: 30,
            timeout: 5,
            failThreshold: 3
        },
        createdAt: new Date(),
        updatedAt: new Date()
    }
]);

// Create example tenant-specific service configurations
// Example: tenant 'acme-corp' uses a custom user service
db.service_configs.insert({
    tenantId: 'acme-corp',
    serviceName: 'user',
    primaryEndpoint: {
        url: 'http://acme-user-service:50052',
        priority: 1,
        weight: 100,
        timeout: 5000,
        headers: {
            'X-Custom-Header': 'acme-value'
        },
        isActive: true
    },
    fallbackChain: [
        {
            url: 'http://acme-user-service-backup:50052',
            priority: 2,
            weight: 50,
            timeout: 5000,
            headers: {},
            isActive: true
        }
    ],
    defaultServiceURL: 'http://go-user-service:50052',
    healthCheck: {
        enabled: true,
        path: '/health',
        method: 'GET',
        interval: 15,
        timeout: 3,
        failThreshold: 3
    },
    loadBalanceStrategy: 'round-robin',
    isActive: true,
    metadata: {
        environment: 'production',
        region: 'us-east-1'
    },
    createdAt: new Date(),
    updatedAt: new Date()
});

print('Migration 002_service_config completed successfully!');
print('Created collections: service_configs, default_service_configs');
print('Created indexes for efficient querying');
print('Inserted 6 default service configurations');
print('Inserted 1 example tenant service configuration');
