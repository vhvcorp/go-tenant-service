# Tenant Service API Documentation

## Base URL
- HTTP: `http://localhost:8083`
- gRPC: `localhost:50053`

## Authentication

All API requests (except health checks) require authentication using one of the following methods:

### API Key Header
```
X-API-Key: your-api-key-here
```

### Bearer Token
```
Authorization: Bearer your-token-here
```

### Tenant Context Headers
```
X-Tenant-ID: tenant-id-here
X-Tenant-Tier: free|basic|professional|enterprise
```

## Rate Limiting

Rate limits are enforced based on subscription tier:
- **Free**: 60 requests/minute
- **Basic**: 300 requests/minute
- **Professional**: 1,000 requests/minute
- **Enterprise**: 5,000 requests/minute

## Endpoints

### Health Check

#### GET /health
Check service health status.

**Response:**
```json
{
  "status": "healthy"
}
```

#### GET /ready
Check if service is ready to accept requests.

**Response:**
```json
{
  "status": "ready"
}
```

---

### Tenant Management

#### POST /api/v1/tenants
Create a new tenant.

**Request Body:**
```json
{
  "name": "Acme Corporation",
  "domain": "acme.com",
  "subscription_tier": "professional"
}
```

**Response:** `201 Created`
```json
{
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Acme Corporation",
    "domain": "acme.com",
    "subscription_tier": "professional",
    "is_active": true,
    "settings": {},
    "created_at": "2025-12-27T10:00:00Z",
    "updated_at": "2025-12-27T10:00:00Z"
  }
}
```

#### GET /api/v1/tenants
List all tenants with pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Items per page (default: 20, max: 100)

**Response:** `200 OK`
```json
{
  "data": {
    "tenants": [
      {
        "id": "507f1f77bcf86cd799439011",
        "name": "Acme Corporation",
        "domain": "acme.com",
        "subscription_tier": "professional",
        "is_active": true,
        "settings": {},
        "created_at": "2025-12-27T10:00:00Z",
        "updated_at": "2025-12-27T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

#### GET /api/v1/tenants/:id
Get tenant details by ID.

**Response:** `200 OK`
```json
{
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Acme Corporation",
    "domain": "acme.com",
    "subscription_tier": "professional",
    "is_active": true,
    "settings": {},
    "created_at": "2025-12-27T10:00:00Z",
    "updated_at": "2025-12-27T10:00:00Z"
  }
}
```

#### PUT /api/v1/tenants/:id
Update tenant information.

**Request Body:**
```json
{
  "name": "Acme Corp",
  "domain": "acmecorp.com",
  "subscription_tier": "enterprise"
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "name": "Acme Corp",
    "domain": "acmecorp.com",
    "subscription_tier": "enterprise",
    "is_active": true,
    "settings": {},
    "created_at": "2025-12-27T10:00:00Z",
    "updated_at": "2025-12-27T10:05:00Z"
  }
}
```

#### DELETE /api/v1/tenants/:id
Delete (soft delete) a tenant.

**Response:** `200 OK`
```json
{
  "message": "Tenant deleted successfully"
}
```

---

### User Management

#### POST /api/v1/tenants/:id/users
Add a user to a tenant.

**Request Body:**
```json
{
  "user_id": "user-123",
  "role": "admin"
}
```

**Response:** `200 OK`
```json
{
  "message": "User added to tenant successfully"
}
```

#### DELETE /api/v1/tenants/:id/users/:user_id
Remove a user from a tenant.

**Response:** `200 OK`
```json
{
  "message": "User removed from tenant successfully"
}
```

---

### Configuration Management

#### GET /api/v1/tenants/:id/configuration
Get tenant configuration settings.

**Response:** `200 OK`
```json
{
  "data": {
    "max_users": 100,
    "features_enabled": ["analytics", "exports"],
    "custom_branding": true
  }
}
```

#### PUT /api/v1/tenants/:id/configuration
Update tenant configuration.

**Request Body:**
```json
{
  "key": "max_users",
  "value": 150,
  "type": "integer"
}
```

**Response:** `200 OK`
```json
{
  "message": "Configuration updated successfully"
}
```

---

### Usage Metrics

#### GET /api/v1/tenants/:id/metrics
Get current usage metrics for a tenant.

**Query Parameters:**
- `period` (optional): Time period (default: "current")

**Response:** `200 OK`
```json
{
  "data": {
    "tenant_id": "507f1f77bcf86cd799439011",
    "api_call_count": 15234,
    "storage_used": 1073741824,
    "bandwidth_used": 5368709120,
    "period": "current",
    "created_at": "2025-12-27T10:00:00Z"
  }
}
```

#### GET /api/v1/tenants/:id/metrics/history
Get historical usage metrics.

**Query Parameters:**
- `limit` (optional): Number of historical records to return (default: 30)

**Response:** `200 OK`
```json
{
  "data": [
    {
      "tenant_id": "507f1f77bcf86cd799439011",
      "api_call_count": 15234,
      "storage_used": 1073741824,
      "bandwidth_used": 5368709120,
      "period": "2025-12-27",
      "created_at": "2025-12-27T10:00:00Z"
    },
    {
      "tenant_id": "507f1f77bcf86cd799439011",
      "api_call_count": 14892,
      "storage_used": 1048576000,
      "bandwidth_used": 5242880000,
      "period": "2025-12-26",
      "created_at": "2025-12-26T10:00:00Z"
    }
  ]
}
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "status_code": 400
  }
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request body or parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

## Usage Examples

### Create a Tenant
```bash
curl -X POST http://localhost:8083/api/v1/tenants \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "name": "Acme Corporation",
    "domain": "acme.com",
    "subscription_tier": "professional"
  }'
```

### Get Tenant Metrics
```bash
curl http://localhost:8083/api/v1/tenants/507f1f77bcf86cd799439011/metrics \
  -H "X-API-Key: your-api-key" \
  -H "X-Tenant-ID: 507f1f77bcf86cd799439011" \
  -H "X-Tenant-Tier: professional"
```

### Update Configuration
```bash
curl -X PUT http://localhost:8083/api/v1/tenants/507f1f77bcf86cd799439011/configuration \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "key": "max_users",
    "value": 200,
    "type": "integer"
  }'
```

---

## Multi-Tenant Isolation

The service supports advanced multi-tenant isolation:

1. **Database Isolation**: Each tenant can have isolated database schemas
2. **Namespace Isolation**: Tenant data is isolated using namespace prefixes
3. **Data Encryption**: Optional at-rest encryption per tenant

Configure isolation settings in the `isolation_config` field when creating or updating tenants.

---

## Subscription Tiers

Available subscription tiers and their limits:

| Tier | Rate Limit | Features |
|------|-----------|----------|
| Free | 60/min | Basic features |
| Basic | 300/min | Standard features + Analytics |
| Professional | 1000/min | Advanced features + Priority support |
| Enterprise | 5000/min | All features + Custom integration + SLA |
