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
- **Semantic tag matching**: Uses multilingual embeddings + LLM suggestions to match tags

### Key Components
- `app/services/scraper.py` - Web scraping with BeautifulSoup
- `app/services/explorer.py` - LLM-driven site exploration
- `app/services/extractor.py` - DSPy extraction for projects and events
- `app/services/tag_matcher.py` - Semantic tag matching with embeddings
- `app/api/routes.py` - FastAPI endpoint (`POST /api/extract`)

### Configuration
Copy `.env.example` to `.env` and set:
- `ANTHROPIC_API_KEY` - Required for LLM calls
- `DSPY_LM_MODEL` - Model to use (default: claude-3-5-haiku-latest)
- `EXPLORER_MAX_PAGES` - Maximum pages to explore (default: 8)
- `EXPLORER_MIN_CONFIDENCE` - Stop exploring when confidence reaches this (default: 0.7)
