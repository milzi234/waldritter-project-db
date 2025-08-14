# Waldritter Website Unified - Development Commands

# Individual Server Commands
# ==========================

# Run Content API server (Go GraphQL server on port 1337)
content-api-dev:
    cd website-content-api && go run main.go -port 1337 -vault ./obsidian-vault

# Run Strapi CMS development server (port 1337) - DEPRECATED
strapi-dev:
    cd website-pages/website-pages && npm run develop

# Run Rails API server (port 3000)
rails-dev:
    cd website-project-db-api && rails server

# Run Admin UI development server
admin-dev:
    cd website-project-db-admin-ui2 && npm run dev

# Run Public Website development server
ui-dev:
    cd website-ui && pnpm dev

# Concurrent Development Commands
# ================================

# Run all services concurrently
dev-all:
    #!/bin/bash
    trap 'kill 0' EXIT
    
    echo "Starting all development servers..."
    echo "=================================="
    echo "Content API: http://localhost:1337"
    echo "Rails API: http://localhost:3000"
    echo "Admin UI: http://localhost:5173"
    echo "Public UI: http://localhost:5174"
    echo "=================================="
    
    (cd website-content-api && go run main.go -port 1337 -vault ./obsidian-vault) &
    (cd website-project-db-api && rails server) &
    (cd website-project-db-admin-ui2 && npm run dev) &
    (cd website-ui && pnpm dev) &
    
    wait

# Run backend services only (Strapi + Rails)
dev-backend:
    #!/bin/bash
    trap 'kill 0' EXIT
    
    echo "Starting backend services..."
    echo "============================"
    echo "Content API: http://localhost:1337"
    echo "Rails API: http://localhost:3000"
    echo "============================"
    
    (cd website-content-api && go run main.go -port 1337 -vault ./obsidian-vault) &
    (cd website-project-db-api && rails server) &
    
    wait

# Run frontend services only (Admin UI + Public Website)
dev-frontend:
    #!/bin/bash
    trap 'kill 0' EXIT
    
    echo "Starting frontend services..."
    echo "============================="
    echo "Admin UI: http://localhost:5173"
    echo "Public UI: http://localhost:5174"
    echo "============================="
    
    (cd website-project-db-admin-ui2 && npm run dev) &
    (cd website-ui && pnpm dev) &
    
    wait

# Utility Commands
# ================

# Install dependencies for all projects
install-all:
    echo "Installing dependencies for all projects..."
    cd website-content-api && go mod download
    cd website-pages/website-pages && npm install
    cd website-project-db-api && bundle install
    cd website-project-db-admin-ui2 && npm install
    cd website-ui && pnpm install
    echo "All dependencies installed!"

# Fix platform-specific binary issues (e.g., esbuild)
fix-platform:
    echo "Fixing platform-specific binary issues..."
    cd website-pages/website-pages && npm rebuild
    cd website-project-db-admin-ui2 && npm rebuild
    echo "Platform binaries rebuilt!"

# Check if servers are running
status:
    #!/bin/bash
    echo "Checking server status..."
    echo ""
    
    if lsof -Pi :1337 -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ Content API is running on port 1337"
    else
        echo "❌ Content API is not running"
    fi
    
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
    
    if lsof -Pi :5174 -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ Public UI is running on port 5174"
    else
        echo "❌ Public UI is not running"
    fi

# Kill processes on development ports if stuck
kill-ports:
    #!/bin/bash
    echo "Killing processes on development ports..."
    
    if lsof -Pi :1337 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:1337)
        echo "Killed process on port 1337"
    fi
    
    if lsof -Pi :3000 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:3000)
        echo "Killed process on port 3000"
    fi
    
    if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:5173)
        echo "Killed process on port 5173"
    fi
    
    if lsof -Pi :5174 -sTCP:LISTEN -t >/dev/null ; then
        kill -9 $(lsof -t -i:5174)
        echo "Killed process on port 5174"
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

# Run Rails console
rails-console:
    cd website-project-db-api && rails console

# Run Rails tests
rails-test:
    cd website-project-db-api && rake test

# Strapi Data Commands
# =====================

# Export Strapi data to a file
strapi-export file="strapi-export.tar.gz":
    cd website-pages/website-pages && npx strapi export --file {{file}} --no-encrypt

# Import Strapi data from a file
strapi-import file="strapi-export.tar.gz":
    cd website-pages/website-pages && npx strapi import --file {{file}} --force

# Export Strapi data with timestamp
strapi-export-backup:
    cd website-pages/website-pages && npx strapi export --file "strapi-backup-$(date +%Y%m%d-%H%M%S).tar.gz" --no-encrypt

# Build Commands
# ==============

# Build Strapi admin panel
strapi-build:
    cd website-pages/website-pages && npm run build

# Build Admin UI for production
admin-build:
    cd website-project-db-admin-ui2 && npm run build

# Build Public Website for production
ui-build:
    cd website-ui && pnpm build

# Build all frontend projects
build-all:
    just strapi-build
    just admin-build
    just ui-build
    echo "All projects built successfully!"

# Testing Commands
# ================

# Run Public Website unit tests
ui-test:
    cd website-ui && pnpm test:unit

# Run Content API tests
content-api-test:
    cd website-content-api && go test ./... -v

# Run all tests
test-all:
    just rails-test
    just ui-test
    just content-api-test
    echo "All tests completed!"

# Content API Commands
# =====================

# Build Content API binary
content-api-build:
    cd website-content-api && go build -o website-content-api main.go

# Run Content API with Docker
content-api-docker:
    cd website-content-api && docker-compose up

# Build Content API Docker image
content-api-docker-build:
    cd website-content-api && docker build -t website-content-api .

# Watch Obsidian vault for changes (verbose mode)
content-api-watch:
    cd website-content-api && go run main.go -port 1337 -vault ./obsidian-vault -verbose

# Open GraphQL playground
content-api-playground:
    open http://localhost:1337

# Docker Commands (for Strapi)
# =============================

# Build Strapi Docker container
strapi-docker-build:
    cd website-pages && docker build -t strapi-toolkit-container .

# Run Strapi in Docker container
strapi-docker-run:
    cd website-pages && docker run -ti --rm -v $(pwd):/project -p 1337:1337 strapi-toolkit-container

# Help Command
# ============

# Show available commands
default:
    @just --list