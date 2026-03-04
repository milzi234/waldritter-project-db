# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo containing three interconnected applications for the Waldritter project management system:

1. **website-project-db-api** - Rails 7 API for project database with GraphQL support
2. **website-project-db-admin-ui2** - Vue 3 admin interface for managing projects
3. **website-url-extractor** - Python FastAPI service for extracting project data from URLs using DSPy

The public website is handled separately by WordPress, which consumes the Rails API for project data. See `docs/WIDGET_SPECIFICATIONS.md` for WordPress widget integration details.

## Development Commands (Justfile)

This project uses [Just](https://just.systems/) for task management. Run `just` to see all available commands.

### Starting Development Servers

```bash
just dev-all        # Start all services (Rails API + Admin UI + URL Extractor)
just dev-backend    # Start Rails API only
just dev-frontend   # Start Admin UI only

# Individual services
just rails-dev      # Start Rails API server (port 3000)
just admin-dev      # Start Admin UI development server (port 5173)
just extractor-dev  # Start URL extractor service (port 8000)
```

### Project Management

```bash
just install-all       # Install dependencies for all projects
just extractor-install # Install Python dependencies for URL extractor
just status            # Check which servers are running
just kill-ports        # Kill stuck processes on development ports
just fix-platform      # Fix platform-specific binary issues
```

### Testing & Building

```bash
just test-all       # Run all tests
just rails-test     # Run Rails tests

just build-all      # Build all frontend projects
just admin-build    # Build Admin UI for production
```

### Database Operations

```bash
just rails-migrate    # Run Rails database migrations
just rails-seed       # Seed Rails database
just rails-setup      # Create, migrate, and seed Rails database
just rails-console    # Open Rails console
```

## Architecture Overview

### Service Communication
- **website-project-db-admin-ui2** connects to Rails API at `http://localhost:3000`
- **WordPress** (external) connects to Rails API at `https://api.waldritter.de/graphql` for project data
- **website-url-extractor** receives proxied requests from Rails API at `http://localhost:8000`
- Authentication handled via Google OAuth in the admin UI (Rails proxies to Python service)

### Key Technologies
- **Frontend**: Vue 3, Vue Router, Pinia (state management), PrimeVue components
- **Backend**: Rails 7 with GraphQL (Project API)
- **URL Extractor**: Python FastAPI, DSPy, BeautifulSoup, sentence-transformers
- **Database**: SQLite for Rails
- **Testing**: Rails built-in test framework, pytest for Python

### Data Models

#### Rails API Core Models
- **Project**: Main content entity with title, description, images
- **Category**: Hierarchical organization for tags
- **Tag**: Labeling system for projects
- **Event**: Recurring events associated with projects
- **Occurrence**: Individual instances of events
- **UmbrellaProject**: Grouping mechanism for related projects

### GraphQL Schema

#### Rails API (`/graphql`)
- Projects (with tags, events, occurrences)
- Categories (hierarchical structure)
- Search functionality
- Event management with recurrence patterns

### Routing Patterns
- **admin-ui**: Protected routes with Google OAuth authentication
- **Rails API**: RESTful + GraphQL endpoints under `/api/v1`

## Important Configuration Files
- `website-project-db-api/config/database.yml` - Rails database config
- `website-project-db-admin-ui2/src/composables/auth_token.js` - Google OAuth authentication
- `website-url-extractor/.env` - Anthropic API key and extractor settings
- `justfile` - Development task definitions
- `docs/WIDGET_SPECIFICATIONS.md` - WordPress widget integration specs

### Storage
- Rails migrations in `website-project-db-api/db/migrate/`
- Active Storage configured for image uploads in Rails

## URL Extractor Service

The URL extractor (`website-url-extractor/`) uses LLM-driven exploration to extract project data from websites:

### Architecture
- **Agentic exploration**: LLM decides which pages to visit based on what information is missing
- **DSPy signatures**: Structured prompts for analysis, link prioritization, and extraction
- **Tag matching**: Keyword matching + LLM suggestions (no embeddings)

### Key Components
- `app/services/scraper.py` - Web scraping with BeautifulSoup
- `app/services/explorer.py` - LLM-driven site exploration
- `app/services/extractor.py` - DSPy extraction for projects and events
- `app/services/tag_matcher.py` - Keyword + LLM tag matching
- `app/api/routes.py` - FastAPI endpoint (`POST /api/extract`)

### Configuration
Copy `.env.example` to `.env` and set:
- `ANTHROPIC_API_KEY` - Required for LLM calls
- `DSPY_LM_MODEL` - Model to use (default: claude-haiku-4-5-20251001)
- `EXPLORER_MAX_PAGES` - Maximum pages to explore (default: 8)
- `EXPLORER_MIN_CONFIDENCE` - Stop exploring when confidence reaches this (default: 0.7)

## Production Deployment

Services are deployed to a k3s cluster on `stauchen.events` via ArgoCD with Kustomize manifests.

### Domain Mapping

| Subdomain | Service | Port |
|---|---|---|
| `project-api.waldritter.dev` | rails-api | 3000 |
| `project-db.waldritter.dev` | admin-ui | 8080 |
| `projects.waldritter.dev` | public-ui (Astro SSR) | 4321 |
| (cluster-internal only) | url-extractor | 8000 |

### Build, Push, Restart Workflow

Images are cross-compiled for `linux/amd64` and transferred via `docker save | scp | ctr import`.

```bash
# Build individual images
just build-api          # Rails API
just build-admin        # Admin UI (bakes in VITE_GOOGLE_CLIENT_ID)
just build-extractor    # URL Extractor
just build-public       # Public UI (Astro SSR, bakes in API_URL)

# Build all
just build-all-images

# Push individual images to server
just push-api
just push-admin
just push-extractor
just push-public

# Push all
just push-all

# Convenience: build + push all
just ship

# Restart a service after pushing (via kubectl over SSH)
ssh root@stauchen.events "KUBECONFIG=/etc/rancher/k3s/k3s.yaml kubectl rollout restart deployment/<name> -n waldritter"
# where <name> is: rails-api, admin-ui, public-ui, url-extractor

# Check deployment status
just prod-status

# View logs
just prod-logs <service>   # e.g. just prod-logs rails-api
```

### Secrets Management

Secrets are stored in `secrets/production/secrets.yaml` (gitignored) and deployed via:
```bash
just deploy-secrets
```

### Key Deployment Notes

- **Admin UI** (`VITE_*` vars): Baked in at Docker build time. Rebuild image to change.
- **Public UI** (`API_URL`): Set as runtime env var in Dockerfile + deployment.yaml. Internal API calls use `http://rails-api:3000`.
- **Rails API**: SQLite on PVC at `/data/db` (not `/app/db` — avoids shadowing migrations). Uses `DATABASE_PATH` env var.
- **url-extractor**: ClusterIP only (no ingress). Rails proxies to it via `EXTRACTOR_SERVICE_URL`.
- **ArgoCD**: Auto-syncs from `deploy/production/` in the GitHub repo. Direct `kubectl apply` changes get reverted — always commit and push, then ArgoCD picks it up.
- **Images**: Must be `--platform linux/amd64` (server is amd64, Mac builds arm64 by default).
- **Registry GC**: Run `just registry-gc` to clean up old image layers.
