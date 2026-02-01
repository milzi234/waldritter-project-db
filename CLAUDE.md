# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo containing two interconnected applications for the Waldritter project management system:

1. **website-project-db-api** - Rails 7 API for project database with GraphQL support
2. **website-project-db-admin-ui2** - Vue 3 admin interface for managing projects

The public website is handled separately by WordPress, which consumes the Rails API for project data. See `docs/WIDGET_SPECIFICATIONS.md` for WordPress widget integration details.

## Development Commands (Justfile)

This project uses [Just](https://just.systems/) for task management. Run `just` to see all available commands.

### Starting Development Servers

```bash
just dev-all        # Start all services (Rails API + Admin UI)
just dev-backend    # Start Rails API only
just dev-frontend   # Start Admin UI only

# Individual services
just rails-dev      # Start Rails API server (port 3000)
just admin-dev      # Start Admin UI development server (port 5173)
```

### Project Management

```bash
just install-all    # Install dependencies for all projects
just status         # Check which servers are running
just kill-ports     # Kill stuck processes on development ports
just fix-platform   # Fix platform-specific binary issues
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
- Authentication handled via Google OAuth in the admin UI

### Key Technologies
- **Frontend**: Vue 3, Vue Router, Pinia (state management), PrimeVue components
- **Backend**: Rails 7 with GraphQL (Project API)
- **Database**: SQLite for Rails
- **Testing**: Rails built-in test framework

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
- `justfile` - Development task definitions
- `docs/WIDGET_SPECIFICATIONS.md` - WordPress widget integration specs

### Storage
- Rails migrations in `website-project-db-api/db/migrate/`
- Active Storage configured for image uploads in Rails
