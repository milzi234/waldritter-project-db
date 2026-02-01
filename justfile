# Waldritter Website Unified - Development Commands

# Individual Server Commands
# ==========================

# Run Rails API server (port 3000)
rails-dev:
    cd website-project-db-api && rails server

# Run Admin UI development server
admin-dev:
    cd website-project-db-admin-ui2 && npm run dev

# Run URL Extractor service (port 8000)
extractor-dev:
    cd website-url-extractor && source venv/bin/activate && uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Concurrent Development Commands
# ================================

# Run all services concurrently
dev-all:
    #!/bin/bash
    trap 'kill 0' EXIT

    echo "Starting all development servers..."
    echo "===================================="
    echo "Rails API: http://localhost:3000"
    echo "Admin UI: http://localhost:5173"
    echo "URL Extractor: http://localhost:8000"
    echo "===================================="

    (cd website-project-db-api && rails server) &
    (cd website-project-db-admin-ui2 && npm run dev) &
    (cd website-url-extractor && source venv/bin/activate && uvicorn app.main:app --reload --host 0.0.0.0 --port 8000) &

    wait

# Run backend services only (Rails API)
dev-backend:
    #!/bin/bash
    trap 'kill 0' EXIT

    echo "Starting backend services..."
    echo "============================"
    echo "Rails API: http://localhost:3000"
    echo "============================"

    cd website-project-db-api && rails server

# Run frontend services only (Admin UI)
dev-frontend:
    #!/bin/bash
    trap 'kill 0' EXIT

    echo "Starting frontend services..."
    echo "============================="
    echo "Admin UI: http://localhost:5173"
    echo "============================="

    cd website-project-db-admin-ui2 && npm run dev

# Utility Commands
# ================

# Install dependencies for all projects
install-all:
    echo "Installing dependencies for all projects..."
    cd website-project-db-api && bundle install
    cd website-project-db-admin-ui2 && npm install
    echo "All dependencies installed!"

# Install dependencies for URL extractor service (uses virtual environment)
extractor-install:
    #!/bin/bash
    echo "Installing URL extractor dependencies..."
    cd website-url-extractor

    # Create virtual environment if it doesn't exist
    if [ ! -d "venv" ]; then
        echo "Creating virtual environment..."
        python3 -m venv venv
    fi

    # Activate and install
    source venv/bin/activate
    pip install -r requirements.txt

    echo "URL extractor dependencies installed in virtual environment!"

# Fix platform-specific binary issues (e.g., esbuild)
fix-platform:
    echo "Fixing platform-specific binary issues..."
    cd website-project-db-admin-ui2 && npm rebuild
    echo "Platform binaries rebuilt!"

# Check if servers are running
status:
    #!/bin/bash
    echo "Checking server status..."
    echo ""

    if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ Rails API is running on port 3000"
    else
        echo "❌ Rails API is not running"
    fi

    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ Admin UI is running on port 5173"
    else
        echo "❌ Admin UI is not running"
    fi

    if lsof -Pi :8000 -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ URL Extractor is running on port 8000"
    else
        echo "❌ URL Extractor is not running"
    fi

# Kill processes on development ports if stuck
kill-ports:
    #!/bin/bash
    echo "Killing processes on development ports..."

    if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:3000)
        echo "Killed process on port 3000"
    fi

    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:5173)
        echo "Killed process on port 5173"
    fi

    if lsof -Pi :8000 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:8000)
        echo "Killed process on port 8000"
    fi

    echo "Done!"

# Rails Database Commands
# ========================

# Run Rails database migrations
rails-migrate:
    cd website-project-db-api && rake db:migrate

# Seed Rails database
rails-seed:
    cd website-project-db-api && rake db:seed

# Setup Rails database (create, migrate, seed)
rails-setup:
    cd website-project-db-api && rake db:create && rake db:migrate && rake db:seed

# Import sample categories and tags data
rails-import-data:
    cd website-project-db-api && ruby import/import.rb

# Run Rails console
rails-console:
    cd website-project-db-api && rails console

# Run Rails tests
rails-test:
    cd website-project-db-api && rake test

# Build Commands
# ==============

# Build Admin UI for production
admin-build:
    cd website-project-db-admin-ui2 && npm run build

# Build all frontend projects
build-all:
    just admin-build
    echo "All projects built successfully!"

# Testing Commands
# ================

# Run all tests
test-all:
    just rails-test
    echo "All tests completed!"

# Kubernetes & k3d Commands
# =========================

# Create k3d cluster with registry
k3d-create:
    #!/bin/bash
    echo "Creating k3d cluster with local registry..."
    k3d cluster create waldritter \
        --api-port 6550 \
        --servers 1 \
        --agents 2 \
        --port "80:80@loadbalancer" \
        --port "443:443@loadbalancer" \
        --registry-create k3d-waldritter-registry.localhost:5001 \
        --volume "$(pwd)/k3d-volumes:/var/lib/rancher/k3s/storage@all" \
        --k3s-arg "--disable=traefik@server:*"

    echo ""
    echo "Waiting for cluster to be ready..."
    kubectl wait --for=condition=ready node --all --timeout=60s

    echo ""
    echo "Cluster created successfully!"
    echo "Registry: k3d-waldritter-registry.localhost:5001"

# Delete k3d cluster
k3d-delete:
    echo "Deleting k3d cluster..."
    k3d cluster delete waldritter

# Start existing k3d cluster
k3d-start:
    echo "Starting k3d cluster..."
    k3d cluster start waldritter

# Stop k3d cluster
k3d-stop:
    echo "Stopping k3d cluster..."
    k3d cluster stop waldritter

# Install nginx ingress controller
k3d-install-ingress:
    #!/bin/bash
    echo "Installing nginx ingress controller..."
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

    echo "Waiting for ingress controller to be ready..."
    kubectl wait --namespace ingress-nginx \
        --for=condition=ready pod \
        --selector=app.kubernetes.io/component=controller \
        --timeout=120s

    echo "Ingress controller installed successfully!"

# Generate local SSL certificates with mkcert
k3d-generate-certs:
    #!/bin/bash
    echo "Generating SSL certificates with mkcert..."

    # Check if mkcert is installed
    if ! command -v mkcert &> /dev/null; then
        echo "mkcert is not installed. Please install it first:"
        echo "  brew install mkcert"
        exit 1
    fi

    # Install local CA
    mkcert -install

    # Create certificates directory
    mkdir -p k3d-certs

    # Generate certificates for all domains
    cd k3d-certs && mkcert \
        "wr-dev.de" \
        "*.wr-dev.de" \
        "admin.wr-dev.de" \
        "api.wr-dev.de"

    # Rename files for clarity
    mv wr-dev.de+3.pem tls.crt
    mv wr-dev.de+3-key.pem tls.key

    echo "Certificates generated in k3d-certs/"
    echo "  - tls.crt"
    echo "  - tls.key"

# Create Kubernetes secret with SSL certificates
k3d-create-tls-secret:
    #!/bin/bash
    echo "Creating TLS secret in Kubernetes..."

    if [ ! -f k3d-certs/tls.crt ] || [ ! -f k3d-certs/tls.key ]; then
        echo "Certificates not found. Run 'just k3d-generate-certs' first."
        exit 1
    fi

    # Create the secret in the waldritter namespace
    kubectl create secret tls waldritter-local-tls \
        --cert=k3d-certs/tls.crt \
        --key=k3d-certs/tls.key \
        --namespace waldritter \
        --dry-run=client -o yaml | kubectl apply -f -

    echo "TLS secret created successfully!"

# Add local domains to /etc/hosts
k3d-setup-hosts:
    #!/bin/bash
    echo "Setting up local domains in /etc/hosts..."
    echo "This requires sudo access."

    # Check if entries already exist
    if grep -q "wr-dev.de" /etc/hosts; then
        echo "Entries already exist in /etc/hosts"
    else
        echo "" | sudo tee -a /etc/hosts
        echo "# Waldritter local development" | sudo tee -a /etc/hosts
        echo "127.0.0.1 wr-dev.de" | sudo tee -a /etc/hosts
        echo "127.0.0.1 admin.wr-dev.de" | sudo tee -a /etc/hosts
        echo "127.0.0.1 api.wr-dev.de" | sudo tee -a /etc/hosts
        echo "Entries added to /etc/hosts"
    fi

# Build all Docker images
k3d-build-images:
    #!/bin/bash
    echo "Building Docker images..."

    # Build Rails API
    echo "Building Rails API..."
    docker build -t k3d-waldritter-registry.localhost:5001/waldritter/rails-api:latest ./website-project-db-api

    # Build Admin UI
    echo "Building Admin UI..."
    docker build -t k3d-waldritter-registry.localhost:5001/waldritter/admin-ui:latest ./website-project-db-admin-ui2

    echo "All images built successfully!"

# Push images to k3d registry
k3d-push-images:
    #!/bin/bash
    echo "Pushing images to k3d registry..."

    docker push k3d-waldritter-registry.localhost:5001/waldritter/rails-api:latest
    docker push k3d-waldritter-registry.localhost:5001/waldritter/admin-ui:latest

    echo "All images pushed successfully!"

# Build and push all images
k3d-build-and-push: k3d-build-images k3d-push-images

# Deploy with Helm
k3d-deploy:
    #!/bin/bash
    echo "Deploying Waldritter website with Helm..."

    helm upgrade --install waldritter ./charts/waldritter-website \
        -f ./charts/waldritter-website/values.local.yaml \
        --create-namespace \
        --namespace waldritter \
        --wait

    echo ""
    echo "Deployment complete!"
    echo "Access the admin interface at:"
    echo "  - https://admin.waldritter.local (Admin interface)"
    echo "  - https://api.waldritter.local (API endpoints)"

# Undeploy Helm release
k3d-undeploy:
    echo "Undeploying Waldritter website..."
    helm uninstall waldritter --namespace waldritter

# Show deployment status
k3d-status:
    #!/bin/bash
    echo "Cluster Status:"
    k3d cluster list
    echo ""
    echo "Pods in waldritter namespace:"
    kubectl get pods -n waldritter
    echo ""
    echo "Services in waldritter namespace:"
    kubectl get svc -n waldritter
    echo ""
    echo "Ingress in waldritter namespace:"
    kubectl get ingress -n waldritter

# View logs for a specific service
k3d-logs service="":
    #!/bin/bash
    if [ -z "{{service}}" ]; then
        echo "Usage: just k3d-logs <service>"
        echo "Available services: rails-api, admin-ui"
        exit 1
    fi
    kubectl logs -n waldritter -l app.kubernetes.io/component={{service}} --tail=100 -f

# Full k3d setup (create cluster, install ingress, generate certs, setup hosts)
k3d-setup: k3d-create k3d-install-ingress k3d-generate-certs k3d-create-tls-secret k3d-setup-hosts
    echo ""
    echo "k3d cluster setup complete!"
    echo "Next steps:"
    echo "  1. Build images: just k3d-build-images"
    echo "  2. Push images: just k3d-push-images"
    echo "  3. Deploy: just k3d-deploy"

# Full k3d teardown
k3d-teardown: k3d-undeploy k3d-delete
    echo "k3d cluster and deployment removed!"

# Quick redeploy (rebuild, push, and deploy)
k3d-redeploy: k3d-build-and-push k3d-deploy
    echo "Application redeployed!"

# Port forward to a specific service (for debugging)
k3d-port-forward service="" port="":
    #!/bin/bash
    if [ -z "{{service}}" ] || [ -z "{{port}}" ]; then
        echo "Usage: just k3d-port-forward <service> <port>"
        echo "Example: just k3d-port-forward rails-api 3000"
        exit 1
    fi
    kubectl port-forward -n waldritter svc/waldritter-{{service}} {{port}}:{{port}}

# Execute into a pod
k3d-exec service="":
    #!/bin/bash
    if [ -z "{{service}}" ]; then
        echo "Usage: just k3d-exec <service>"
        echo "Available services: rails-api, admin-ui"
        exit 1
    fi
    POD=$(kubectl get pods -n waldritter -l app.kubernetes.io/component={{service}} -o jsonpath='{.items[0].metadata.name}')
    kubectl exec -it -n waldritter $POD -- /bin/sh

# Help Command
# ============

# Show available commands
default:
    @just --list
