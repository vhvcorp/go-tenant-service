# Tenant Service Architecture

## Overview

The Tenant Service is a microservice responsible for managing multi-tenant organizations, their configurations, users, and usage metrics within the SaaS framework. It provides both HTTP/REST and gRPC interfaces for maximum flexibility.

## Architecture Layers

### 1. Presentation Layer (HTTP/gRPC)
- **HTTP Server (Gin)**: RESTful API endpoints on port 8083
- **gRPC Server**: High-performance RPC on port 50053
- **Middleware Stack**:
  - Recovery middleware for panic handling
  - Tenant context extraction
  - API key authentication
  - Rate limiting (tier-based)
  - Usage tracking

### 2. Application Layer (Handlers/Services)
- **Handlers**: Process HTTP requests and format responses
- **Services**: Business logic and orchestration
- **Request/Response Models**: Data transfer objects

### 3. Domain Layer
- **Models**: Core domain entities (Tenant, TenantUser, UsageMetrics)
- **Business Rules**: Validation and business constraints
- **Interfaces**: Repository contracts

### 4. Infrastructure Layer
- **Repositories**: Data access implementations
- **MongoDB**: Primary data store
- **Indexes**: Optimized for query patterns

## Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     External Clients                         │
│          (Web Apps, Mobile Apps, Other Services)             │
└───────────────────┬─────────────────────────────────────────┘
                    │
                    ├──────────┬──────────┐
                    │          │          │
            ┌───────▼─────┐    │    ┌─────▼──────┐
            │ HTTP:8083   │    │    │ gRPC:50053 │
            │ (REST API)  │    │    │            │
            └───────┬─────┘    │    └─────┬──────┘
                    │          │          │
            ┌───────▼──────────▼──────────▼───────┐
            │         Middleware Stack            │
            │  ┌──────────────────────────────┐   │
            │  │ - Tenant Context             │   │
            │  │ - Authentication             │   │
            │  │ - Rate Limiting              │   │
            │  │ - Usage Tracking             │   │
            │  └──────────────────────────────┘   │
            └───────────────┬─────────────────────┘
                            │
            ┌───────────────▼─────────────────────┐
            │          Handlers Layer             │
            │  ┌──────────────────────────────┐   │
            │  │ TenantHandler                │   │
            │  │ - CRUD operations            │   │
            │  │ - Configuration management   │   │
            │  │ - Usage metrics              │   │
            │  └──────────────────────────────┘   │
            └───────────────┬─────────────────────┘
                            │
            ┌───────────────▼─────────────────────┐
            │          Services Layer             │
            │  ┌──────────────────────────────┐   │
            │  │ TenantService                │   │
            │  │ - Business logic             │   │
            │  │ - Validation                 │   │
            │  │ - Orchestration              │   │
            │  └──────────────────────────────┘   │
            └───────────────┬─────────────────────┘
                            │
            ┌───────────────▼─────────────────────┐
            │        Repositories Layer           │
            │  ┌──────────────────────────────┐   │
            │  │ - TenantRepository           │   │
            │  │ - TenantUserRepository       │   │
            │  │ - UsageMetricsRepository     │   │
            │  └──────────────────────────────┘   │
            └───────────────┬─────────────────────┘
                            │
            ┌───────────────▼─────────────────────┐
            │          MongoDB Database           │
            │  ┌──────────────────────────────┐   │
            │  │ Collections:                 │   │
            │  │ - tenants                    │   │
            │  │ - tenant_users               │   │
            │  │ - usage_metrics              │   │
            │  └──────────────────────────────┘   │
            └─────────────────────────────────────┘
```

## Multi-Tenant Isolation Strategy

### 1. Database-Level Isolation
- **Shared Database, Shared Schema**: Default mode with tenant_id filtering
- **Shared Database, Separate Schemas**: Optional namespace-based isolation
- **Separate Databases**: Enterprise tier option (configured per tenant)

### 2. Application-Level Isolation
- Tenant context middleware ensures all requests include tenant identification
- Repository layer automatically filters by tenant_id
- No cross-tenant data access without explicit permission

### 3. Security Isolation
- API keys scoped to specific tenants
- Rate limiting per tenant based on subscription tier
- Optional data encryption at rest per tenant

## Data Model

### Core Entities

#### Tenant
```go
{
  ID               ObjectID
  Name             string
  Domain           string
  SubscriptionTier string
  IsActive         bool
  Settings         map[string]interface{}
  IsolationConfig  *IsolationConfig
  CreatedAt        time.Time
  UpdatedAt        time.Time
}
```

#### IsolationConfig
```go
{
  DatabaseIsolation bool
  Namespace         string
  DataEncryption    bool
}
```

#### TenantUser
```go
{
  ID        ObjectID
  TenantID  string
  UserID    string
  Role      string
  IsActive  bool
  CreatedAt time.Time
  UpdatedAt time.Time
}
```

#### UsageMetrics
```go
{
  ID            ObjectID
  TenantID      string
  APICallCount  int64
  StorageUsed   int64
  BandwidthUsed int64
  Period        string
  CreatedAt     time.Time
}
```

## Security Architecture

### 1. Authentication
- **API Key Authentication**: Validates X-API-Key header
- **Bearer Token**: Supports OAuth 2.0 compatible tokens
- **Service-to-Service**: mTLS for gRPC communication (future)

### 2. Authorization
- Tenant context validation
- Role-based access control (RBAC) via tenant-user relationships
- Resource ownership verification

### 3. Rate Limiting
Tier-based token bucket algorithm:
- Tokens refill per minute based on subscription tier
- Concurrent request tracking
- Graceful degradation with 429 responses

### 4. Data Protection
- Sensitive data encryption in transit (TLS)
- Optional at-rest encryption per tenant
- Audit logging for security events

## Scalability Considerations

### Horizontal Scaling
- Stateless service design
- No in-memory session state
- MongoDB connection pooling
- Load balancer ready

### Vertical Scaling
- Configurable MongoDB pool sizes
- Optimized database indexes
- Efficient query patterns

### Performance Optimizations
- Database indexes on tenant_id, name, domain
- Pagination for list operations
- TTL indexes for usage metrics (30-day retention)
- Aggregation pipelines for metrics

## Monitoring & Observability

### Metrics
- API call count per tenant
- Request latency per endpoint
- Error rates
- Resource usage (storage, bandwidth)

### Logging
- Structured logging with Zap
- Log levels: debug, info, warn, error
- Contextual logging with tenant_id

### Health Checks
- `/health`: Basic service health
- `/ready`: Readiness check with dependencies

## Deployment Architecture

### Container-Based Deployment
```
┌─────────────────────────────────────────────────┐
│              Load Balancer                      │
└──────────┬──────────────┬───────────────────────┘
           │              │
    ┌──────▼─────┐  ┌─────▼──────┐
    │ Instance 1 │  │ Instance 2 │  (Auto-scaling)
    └──────┬─────┘  └─────┬──────┘
           │              │
    ┌──────▼──────────────▼──────┐
    │    MongoDB Cluster          │
    │  (Replica Set or Sharded)   │
    └─────────────────────────────┘
```

### Environment Configuration
- Development: Single instance, local MongoDB
- Staging: 2 instances, shared MongoDB
- Production: Auto-scaling (2-10 instances), MongoDB cluster

## Inter-Service Communication

### Services that call Tenant Service:
- User Service: Tenant validation
- All Services: Tenant context verification

### Services called by Tenant Service:
- None (leaf service in dependency graph)

## Future Enhancements

1. **Caching Layer**
   - Redis for tenant metadata
   - Cache invalidation strategies

2. **Event-Driven Architecture**
   - Publish tenant events (created, updated, deleted)
   - RabbitMQ/Kafka integration

3. **Advanced Analytics**
   - Real-time dashboards
   - Predictive analytics for resource usage

4. **Multi-Region Support**
   - Geographic data distribution
   - Region-based routing

5. **Backup and Recovery**
   - Automated backups per tenant
   - Point-in-time recovery

## Technology Stack

- **Language**: Go 1.25.5
- **Web Framework**: Gin
- **RPC Framework**: gRPC
- **Database**: MongoDB 4.4+
- **Logging**: Uber Zap
- **Configuration**: Viper
- **Containerization**: Docker
- **Orchestration**: Kubernetes (recommended)

## Development Guidelines

### Code Organization
```
cmd/              - Application entry points
internal/
  domain/         - Domain models and interfaces
  handler/        - HTTP handlers
  service/        - Business logic
  repository/     - Data access layer
  middleware/     - HTTP middleware
  grpc/           - gRPC server implementation
```

### Best Practices
1. Always use context for cancellation and timeouts
2. Log with structured fields (tenant_id, user_id, etc.)
3. Handle errors explicitly, don't panic
4. Use interfaces for testability
5. Follow Go naming conventions
6. Document public APIs

### Testing Strategy
1. Unit tests for business logic
2. Integration tests for repositories
3. End-to-end tests for critical flows
4. Load testing for performance validation
