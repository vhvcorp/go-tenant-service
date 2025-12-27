# Tenant Service Deployment Guide

## Prerequisites

- Docker 20.10+
- Kubernetes 1.21+ (for production)
- MongoDB 4.4+
- Go 1.25.5+ (for local development)

## Environment Variables

See [DEPENDENCIES.md](DEPENDENCIES.md) for complete list of environment variables.

Required variables:
```bash
TENANT_SERVICE_PORT=50053
TENANT_SERVICE_HTTP_PORT=8083
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=saas_framework
LOG_LEVEL=info
```

## Local Development Deployment

### 1. Using Go directly

```bash
# Install dependencies
make deps

# Run the service
make run
```

### 2. Using Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

### 3. Using Docker Compose (recommended for local dev)

Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    environment:
      MONGO_INITDB_DATABASE: saas_framework

  tenant-service:
    build: .
    ports:
      - "8083:8083"
      - "50053:50053"
    environment:
      MONGODB_URI: mongodb://mongodb:27017
      MONGODB_DATABASE: saas_framework
      LOG_LEVEL: debug
    depends_on:
      - mongodb

volumes:
  mongodb_data:
```

Run with:
```bash
docker-compose up -d
```

## Production Deployment

### Option 1: Kubernetes Deployment

#### 1. Build and push Docker image

```bash
# Set your registry
export DOCKER_REGISTRY=ghcr.io/vhvplatform

# Build and push
make docker-build
make docker-push
```

#### 2. Create Kubernetes resources

**ConfigMap** (`k8s/configmap.yaml`):
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tenant-service-config
  namespace: production
data:
  TENANT_SERVICE_PORT: "50053"
  TENANT_SERVICE_HTTP_PORT: "8083"
  MONGODB_DATABASE: "saas_framework"
  LOG_LEVEL: "info"
  ENVIRONMENT: "production"
```

**Secret** (`k8s/secret.yaml`):
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tenant-service-secret
  namespace: production
type: Opaque
stringData:
  MONGODB_URI: "mongodb://username:password@mongodb-service:27017"
```

**Deployment** (`k8s/deployment.yaml`):
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tenant-service
  namespace: production
  labels:
    app: tenant-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: tenant-service
  template:
    metadata:
      labels:
        app: tenant-service
    spec:
      containers:
      - name: tenant-service
        image: ghcr.io/vhvplatform/tenant-service:latest
        ports:
        - containerPort: 8083
          name: http
          protocol: TCP
        - containerPort: 50053
          name: grpc
          protocol: TCP
        envFrom:
        - configMapRef:
            name: tenant-service-config
        - secretRef:
            name: tenant-service-secret
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8083
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8083
          initialDelaySeconds: 5
          periodSeconds: 5
```

**Service** (`k8s/service.yaml`):
```yaml
apiVersion: v1
kind: Service
metadata:
  name: tenant-service
  namespace: production
spec:
  selector:
    app: tenant-service
  ports:
  - name: http
    port: 8083
    targetPort: 8083
    protocol: TCP
  - name: grpc
    port: 50053
    targetPort: 50053
    protocol: TCP
  type: ClusterIP
```

**HorizontalPodAutoscaler** (`k8s/hpa.yaml`):
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: tenant-service-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: tenant-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

#### 3. Apply Kubernetes resources

```bash
# Create namespace
kubectl create namespace production

# Apply resources
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/hpa.yaml

# Verify deployment
kubectl get pods -n production
kubectl get svc -n production
```

### Option 2: Cloud Platform Deployment

#### AWS ECS

1. Create ECR repository:
```bash
aws ecr create-repository --repository-name tenant-service
```

2. Build and push:
```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker build -t tenant-service .
docker tag tenant-service:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/tenant-service:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/tenant-service:latest
```

3. Create ECS task definition and service using AWS Console or CLI

#### Google Cloud Run

```bash
# Build and push to GCR
gcloud builds submit --tag gcr.io/PROJECT_ID/tenant-service

# Deploy
gcloud run deploy tenant-service \
  --image gcr.io/PROJECT_ID/tenant-service \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8083
```

#### Azure Container Instances

```bash
# Login to Azure
az login

# Create container registry
az acr create --resource-group myResourceGroup --name myregistry --sku Basic

# Build and push
az acr build --registry myregistry --image tenant-service:latest .

# Deploy
az container create \
  --resource-group myResourceGroup \
  --name tenant-service \
  --image myregistry.azurecr.io/tenant-service:latest \
  --cpu 1 --memory 1 \
  --ports 8083 50053
```

## Tenant-Specific Deployments

For enterprise customers requiring dedicated instances:

```bash
# Build tenant-specific image
make deploy-tenant TENANT_ID=acme-corp

# Or with Docker directly
docker build --build-arg TENANT_ID=acme-corp \
  -t tenant-service:acme-corp .
```

Deploy using the same Kubernetes manifests with updated:
- Namespace: `production-acme-corp`
- Image tag: `tenant-service:acme-corp`
- Resource quotas as per SLA

## Database Setup

### MongoDB Setup

#### Development (Single Instance)
```bash
docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -v mongodb_data:/data/db \
  mongo:7.0
```

#### Production (Replica Set)

Use MongoDB Atlas or self-hosted replica set:

1. **MongoDB Atlas** (Recommended):
   - Create cluster at https://cloud.mongodb.com
   - Configure network access
   - Get connection string
   - Update `MONGODB_URI` environment variable

2. **Self-Hosted Replica Set**:
```bash
# Initialize replica set
docker run -d --name mongo1 -p 27017:27017 mongo:7.0 --replSet rs0
docker run -d --name mongo2 -p 27018:27017 mongo:7.0 --replSet rs0
docker run -d --name mongo3 -p 27019:27017 mongo:7.0 --replSet rs0

# Configure replica set
docker exec -it mongo1 mongosh
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongo1:27017" },
    { _id: 1, host: "mongo2:27017" },
    { _id: 2, host: "mongo3:27017" }
  ]
})
```

### Database Indexes

Indexes are created automatically on service startup, but you can manually create them:

```javascript
// Connect to MongoDB
use saas_framework;

// Tenants collection
db.tenants.createIndex({ "name": 1 }, { unique: true });
db.tenants.createIndex({ "domain": 1 }, { unique: true, sparse: true });

// Tenant users collection
db.tenant_users.createIndex({ "tenant_id": 1, "user_id": 1 }, { unique: true });

// Usage metrics collection
db.usage_metrics.createIndex({ "tenant_id": 1, "period": 1 });
db.usage_metrics.createIndex({ "created_at": -1 }, { expireAfterSeconds: 2592000 }); // 30 days TTL
```

## Monitoring Setup

### Prometheus Metrics (Future Enhancement)

Add Prometheus annotations to Kubernetes deployment:
```yaml
metadata:
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8083"
    prometheus.io/path: "/metrics"
```

### Logging

Logs are written to stdout/stderr in JSON format. Configure log aggregation:

**Fluentd (Kubernetes)**:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/tenant-service*.log
      pos_file /var/log/tenant-service.log.pos
      tag tenant-service
      format json
    </source>
```

## Health Checks

The service provides two health check endpoints:

- `/health`: Basic health check (always returns 200 if service is running)
- `/ready`: Readiness check (verifies database connectivity)

## Rollback Procedure

### Kubernetes
```bash
# View deployment history
kubectl rollout history deployment/tenant-service -n production

# Rollback to previous version
kubectl rollout undo deployment/tenant-service -n production

# Rollback to specific revision
kubectl rollout undo deployment/tenant-service -n production --to-revision=2
```

### Docker
```bash
# Tag and push previous version as latest
docker tag tenant-service:v1.2.3 tenant-service:latest
docker push tenant-service:latest

# Restart containers
docker-compose restart tenant-service
```

## Performance Tuning

### MongoDB Connection Pool
Adjust based on load:
```bash
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=10
```

### Resource Limits
Adjust Kubernetes resource limits based on observed usage:
```yaml
resources:
  requests:
    memory: "1Gi"
    cpu: "1000m"
  limits:
    memory: "2Gi"
    cpu: "2000m"
```

## Troubleshooting

### Service won't start
1. Check logs: `kubectl logs -f deployment/tenant-service -n production`
2. Verify MongoDB connectivity
3. Check environment variables

### High CPU usage
1. Review API call patterns
2. Check for missing database indexes
3. Enable query profiling in MongoDB

### Memory leaks
1. Monitor with `pprof` (add endpoint in future)
2. Check for goroutine leaks
3. Review connection pool settings

## Security Hardening

1. **Enable TLS**:
   - Use cert-manager in Kubernetes
   - Mount certificates in container

2. **Network Policies**:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: tenant-service-policy
spec:
  podSelector:
    matchLabels:
      app: tenant-service
  ingress:
  - from:
    - podSelector:
        matchLabels:
          role: api-gateway
```

3. **Run security scans**:
```bash
make security-scan
```

## Backup and Recovery

### Database Backups

**Automated backups** (MongoDB Atlas):
- Configure backup schedule in Atlas UI
- Default: Point-in-time recovery with 7-day retention

**Manual backup**:
```bash
mongodump --uri="mongodb://localhost:27017/saas_framework" --out=/backup
```

**Restore**:
```bash
mongorestore --uri="mongodb://localhost:27017/saas_framework" /backup/saas_framework
```

## Support

For deployment issues:
- Check logs: `kubectl logs -f deployment/tenant-service`
- Review [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- See [API.md](API.md) for endpoint documentation
