# Tenant Service Dependencies

## Shared Packages (from go-shared)

```go
require (
    github.com/vhvcorp/go-shared/config
    github.com/vhvcorp/go-shared/logger
    github.com/vhvcorp/go-shared/mongodb
    github.com/vhvcorp/go-shared/redis
    github.com/vhvcorp/go-shared/errors
    github.com/vhvcorp/go-shared/middleware
    github.com/vhvcorp/go-shared/response
    github.com/vhvcorp/go-shared/validation
    github.com/vhvcorp/go-shared/tenant
)
```

## External Dependencies

### Infrastructure
- **MongoDB**: Tenant data, subscriptions, domains
  - Collections: `tenants`, `subscriptions`, `domains`
- **Redis**: Tenant cache
  - Keys: `tenant:*`, `subscription:*`

### Third-party Libraries
```go
require (
    github.com/gin-gonic/gin v1.10.0
    google.golang.org/grpc v1.69.2
    go.mongodb.org/mongo-driver v1.17.3
)
```

## Inter-service Communication

### Services Called by Tenant Service
- None (leaf service)

### Services Calling Tenant Service
- **User Service**: Tenant validation
- **All Services**: Tenant context verification

## Environment Variables

```bash
# Server
TENANT_SERVICE_PORT=50053
TENANT_SERVICE_HTTP_PORT=8083

# Database
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=saas_framework

# Redis
REDIS_URL=redis://localhost:6379/0

# Logging
LOG_LEVEL=info
```

## Database Schema

### Collections

#### tenants
```json
{
  "_id": "ObjectId",
  "name": "string (indexed)",
  "slug": "string (unique, indexed)",
  "subscription_tier": "string",
  "status": "string (indexed)",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### domains
```json
{
  "_id": "ObjectId",
  "tenant_id": "string (indexed)",
  "domain": "string (unique, indexed)",
  "is_verified": "boolean",
  "created_at": "timestamp"
}
```

## Resource Requirements

### Production
- CPU: 1 core
- Memory: 1GB
- Replicas: 2
