# Waldritter Website - Kubernetes Deployment Guide

This guide provides comprehensive instructions for deploying the Waldritter website application suite to Kubernetes using Docker, Helm, and k3d for local development or production environments.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Prerequisites](#prerequisites)
3. [Local Development with k3d](#local-development-with-k3d)
4. [Production Deployment](#production-deployment)
5. [Configuration](#configuration)
6. [Troubleshooting](#troubleshooting)
7. [Maintenance](#maintenance)

## Architecture Overview

The Waldritter website consists of four interconnected services:

| Service | Description | Port | Technology |
|---------|-------------|------|------------|
| **content-api** | GraphQL API serving Obsidian markdown content | 1337 | Go |
| **rails-api** | Rails API for project database | 3000 | Ruby on Rails 7 |
| **website-ui** | Public-facing website | 80 | Vue 3 + Nginx |
| **admin-ui** | Admin interface with Google OAuth | 80 | Vue 3 + Nginx |

### Data Flow

```
┌─────────────────┐     ┌─────────────────┐
│   website-ui    │────▶│   content-api   │
│   (Vue + Nginx) │     │  (Go GraphQL)   │
└────────┬────────┘     └─────────────────┘
         │
         │              ┌─────────────────┐
         └─────────────▶│    rails-api    │
                        │  (Rails + SQLite)│
┌─────────────────┐     └─────────────────┘
│    admin-ui     │────────────┘
│ (Vue + Nginx)   │
└─────────────────┘
```

## Prerequisites

### Required Tools

- **Docker** (20.10+): Container runtime
- **kubectl** (1.25+): Kubernetes CLI
- **Helm** (3.10+): Kubernetes package manager
- **k3d** (5.0+): Local Kubernetes clusters
- **mkcert**: Local SSL certificates
- **Helmfile** (optional): Multi-environment deployments

### Installation (macOS)

```bash
# Install Docker Desktop
brew install --cask docker

# Install Kubernetes tools
brew install kubectl helm k3d mkcert

# Install Helmfile (optional)
brew install helmfile

# Install Just (task runner)
brew install just
```

### Installation (Linux)

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Install k3d
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install mkcert
brew install mkcert  # or download from GitHub releases

# Install Just
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin
```

## Local Development with k3d

### Quick Start

```bash
# 1. Complete setup (cluster, ingress, SSL, hosts)
just k3d-setup

# 2. Build and push Docker images
just k3d-build-and-push

# 3. Deploy the application
just k3d-deploy

# Access the application
# - https://waldritter.local (Main website)
# - https://admin.waldritter.local (Admin interface)
# - https://api.waldritter.local (API endpoints)
```

### Step-by-Step Setup

#### 1. Create k3d Cluster

```bash
just k3d-create
```

This creates a k3d cluster with:
- Local Docker registry at `k3d-waldritter-registry.localhost:5000`
- LoadBalancer ports 80 and 443
- Persistent volumes
- Traefik disabled (using nginx-ingress instead)

#### 2. Install Ingress Controller

```bash
just k3d-install-ingress
```

Installs the nginx-ingress controller for routing traffic.

#### 3. Generate SSL Certificates

```bash
just k3d-generate-certs
```

Uses mkcert to generate trusted SSL certificates for:
- `waldritter.local`
- `*.waldritter.local`
- `admin.waldritter.local`
- `api.waldritter.local`

#### 4. Create TLS Secret

```bash
just k3d-create-tls-secret
```

Creates a Kubernetes secret with the SSL certificates.

#### 5. Setup Local DNS

```bash
just k3d-setup-hosts
```

Adds entries to `/etc/hosts` for local domain resolution.

#### 6. Build Docker Images

```bash
just k3d-build-images
```

Builds Docker images for all services:
- `content-api`: Go GraphQL server
- `rails-api`: Rails API with SQLite
- `website-ui`: Vue 3 public website
- `admin-ui`: Vue 3 admin interface

#### 7. Push to Registry

```bash
just k3d-push-images
```

Pushes images to the local k3d registry.

#### 8. Deploy Application

```bash
just k3d-deploy
```

Deploys the application using Helm with local values.

### Managing the Deployment

```bash
# Check status
just k3d-status

# View logs
just k3d-logs content-api
just k3d-logs rails-api
just k3d-logs website-ui
just k3d-logs admin-ui

# Execute into a pod
just k3d-exec rails-api

# Port forward for debugging
just k3d-port-forward rails-api 3000

# Redeploy after changes
just k3d-redeploy

# Remove deployment
just k3d-undeploy

# Delete cluster
just k3d-delete
```

## Production Deployment

### Prerequisites

1. **Kubernetes Cluster**: EKS, GKE, AKS, or self-managed
2. **Domain Names**: Configured DNS pointing to cluster
3. **cert-manager**: For Let's Encrypt SSL certificates
4. **Container Registry**: Docker Hub, ECR, GCR, etc.

### Setup

#### 1. Install cert-manager

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml
```

#### 2. Configure Secrets

Create production secrets:

```bash
cp secrets/production/secrets.yaml.example secrets/production/secrets.yaml
# Edit secrets.yaml with actual values

# Generate Rails secret
cd website-project-db-api
rails secret
```

#### 3. Build and Push Images

```bash
# Build production images
docker build -t your-registry/waldritter/content-api:v1.0.0 ./website-content-api
docker build -t your-registry/waldritter/rails-api:v1.0.0 ./website-project-db-api
docker build -t your-registry/waldritter/website-ui:v1.0.0 ./website-ui
docker build -t your-registry/waldritter/admin-ui:v1.0.0 ./website-project-db-admin-ui2

# Push to registry
docker push your-registry/waldritter/content-api:v1.0.0
docker push your-registry/waldritter/rails-api:v1.0.0
docker push your-registry/waldritter/website-ui:v1.0.0
docker push your-registry/waldritter/admin-ui:v1.0.0
```

#### 4. Deploy with Helmfile

```bash
# Deploy to production
helmfile -e production sync

# Or deploy with Helm directly
helm upgrade --install waldritter ./charts/waldritter-website \
  -f ./charts/waldritter-website/values.production.yaml \
  --namespace waldritter \
  --create-namespace
```

### Production Checklist

- [ ] SSL certificates configured (cert-manager)
- [ ] Secrets properly configured
- [ ] Database backups configured
- [ ] Monitoring/logging setup
- [ ] Resource limits appropriate
- [ ] Autoscaling configured
- [ ] Network policies enabled
- [ ] RBAC configured
- [ ] Ingress annotations for rate limiting
- [ ] Health checks configured

## Configuration

### Environment Variables

#### Content API
- `PORT`: Server port (default: 1337)
- `PLAYGROUND_ENABLED`: Enable GraphQL playground
- `WATCH_ENABLED`: Watch for file changes

#### Rails API
- `RAILS_ENV`: Environment (development/production)
- `SECRET_KEY_BASE`: Rails secret key
- `RAILS_LOG_TO_STDOUT`: Enable stdout logging
- `SEED_DATABASE`: Seed database on startup

#### Website UI
- `VITE_GRAPHQL_URI`: Content API GraphQL endpoint
- `VITE_RAILS_API_URL`: Rails API URL
- `VITE_APP_ENV`: Environment

#### Admin UI
- `VITE_API_BASE_URL`: Rails API URL
- `VITE_GOOGLE_CLIENT_ID`: Google OAuth client ID
- `VITE_APP_ENV`: Environment

### Helm Values

Key configuration in `values.yaml`:

```yaml
# Image configuration
contentApi:
  image:
    repository: waldritter/content-api
    tag: latest
  replicaCount: 2
  resources:
    requests:
      cpu: 100m
      memory: 128Mi

# Persistence
railsApi:
  persistence:
    database:
      enabled: true
      size: 5Gi
    storage:
      enabled: true
      size: 10Gi

# Autoscaling
websiteUi:
  autoscaling:
    enabled: true
    minReplicas: 2
    maxReplicas: 10

# TLS configuration
ingress:
  tls:
    provider: cert-manager  # or mkcert for local
    certManager:
      issuer: letsencrypt-prod
      email: admin@waldritter.com
```

## Troubleshooting

### Common Issues

#### 1. Pods Not Starting

```bash
# Check pod status
kubectl get pods -n waldritter

# Check pod logs
kubectl logs -n waldritter <pod-name>

# Describe pod for events
kubectl describe pod -n waldritter <pod-name>
```

#### 2. Image Pull Errors

```bash
# Check if images exist in registry
docker images | grep waldritter

# For k3d, ensure registry is running
docker ps | grep registry

# Re-push images
just k3d-push-images
```

#### 3. SSL Certificate Issues

```bash
# Check certificate status
kubectl get certificates -n waldritter

# Check cert-manager logs
kubectl logs -n cert-manager deploy/cert-manager

# For local, regenerate certificates
just k3d-generate-certs
just k3d-create-tls-secret
```

#### 4. Database Connection Issues

```bash
# Check Rails API logs
kubectl logs -n waldritter -l app.kubernetes.io/component=rails-api

# Execute into pod and test database
kubectl exec -it -n waldritter <rails-pod> -- rails console
```

#### 5. Ingress Not Working

```bash
# Check ingress status
kubectl get ingress -n waldritter

# Check nginx-ingress logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx

# Test with port-forward
kubectl port-forward -n waldritter svc/waldritter-website-ui 8080:80
```

### Debug Commands

```bash
# Get all resources
kubectl get all -n waldritter

# Check events
kubectl get events -n waldritter --sort-by='.lastTimestamp'

# Resource usage
kubectl top pods -n waldritter

# Check persistent volumes
kubectl get pv,pvc -n waldritter

# Test connectivity
kubectl run -it --rm debug --image=alpine --restart=Never -n waldritter -- sh
# Inside pod:
# apk add curl
# curl http://waldritter-rails-api:3000/health
```

## Maintenance

### Backup Procedures

#### Database Backup

```bash
# Create backup
kubectl exec -n waldritter <rails-pod> -- \
  sqlite3 /app/db/production.sqlite3 ".backup /tmp/backup.db"

# Copy backup locally
kubectl cp waldritter/<rails-pod>:/tmp/backup.db ./backup-$(date +%Y%m%d).db
```

#### Obsidian Vault Backup

```bash
# Copy vault content
kubectl cp waldritter/<content-api-pod>:/root/obsidian-vault ./vault-backup-$(date +%Y%m%d)
```

### Updating Applications

#### Rolling Update

```bash
# Update image tag
kubectl set image deployment/waldritter-website-ui \
  website-ui=waldritter/website-ui:v1.0.1 \
  -n waldritter

# Or use Helm
helm upgrade waldritter ./charts/waldritter-website \
  --set websiteUi.image.tag=v1.0.1 \
  -n waldritter
```

#### Database Migrations

```bash
# Run migrations
kubectl exec -n waldritter <rails-pod> -- rails db:migrate

# Or create a Job
kubectl create job rails-migrate \
  --from=deployment/waldritter-rails-api \
  -n waldritter -- rails db:migrate
```

### Monitoring

#### Health Checks

All services expose `/health` endpoints:

```bash
# Check health endpoints
curl https://waldritter.local/health
curl https://api.waldritter.local/health
curl https://admin.waldritter.local/health
```

#### Metrics

Configure Prometheus/Grafana for monitoring:

```yaml
# Example ServiceMonitor for Prometheus Operator
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: waldritter
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: waldritter-website
  endpoints:
  - port: http
    path: /metrics
```

### Scaling

#### Manual Scaling

```bash
# Scale deployment
kubectl scale deployment waldritter-website-ui --replicas=5 -n waldritter
```

#### Autoscaling

HPA is configured for website-ui in production:

```bash
# Check HPA status
kubectl get hpa -n waldritter

# Modify HPA
kubectl edit hpa waldritter-website-ui -n waldritter
```

## Security Considerations

1. **Secrets Management**: Use external secret stores (Vault, AWS Secrets Manager)
2. **Network Policies**: Enable in production to restrict pod communication
3. **RBAC**: Use ServiceAccounts with minimal permissions
4. **Image Scanning**: Scan images for vulnerabilities before deployment
5. **Pod Security Standards**: Apply restricted security contexts
6. **Ingress Security**: Enable rate limiting and WAF rules

## Support

For issues or questions:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review application logs with `just k3d-logs <service>`
3. Create an issue in the repository with:
   - Error messages
   - Steps to reproduce
   - Environment details
   - Relevant configuration