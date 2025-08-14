# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo containing four interconnected web applications for a Waldritter website system:

1. **website-content-api** - Go-based GraphQL content API serving Obsidian markdown files with natural markdown syntax
2. **website-project-db-api** - Rails 7 API for project database with GraphQL support  
3. **website-project-db-admin-ui2** - Vue 3 admin interface for managing projects
4. **website-ui** - Vue 3 public-facing website with Apollo GraphQL client

## Development Commands (Justfile)

This project uses [Just](https://just.systems/) for task management. Run `just` to see all available commands.

### Starting Development Servers

```bash
just dev-all        # Start all services (Content API, Rails API, Admin UI, Public Website)
just dev-backend    # Start backend services only (Content API + Rails API)
just dev-frontend   # Start frontend services only (Admin UI + Public Website)

# Individual services
just content-api-dev    # Start Content API (Go GraphQL server on port 1337)
just rails-dev         # Start Rails API server (port 3000)
just admin-dev         # Start Admin UI development server
just ui-dev            # Start Public Website development server
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
just test-all          # Run all tests
just rails-test        # Run Rails tests
just ui-test          # Run Vue unit tests
just content-api-test # Run Go tests

just build-all        # Build all frontend projects
just admin-build      # Build Admin UI for production
just ui-build         # Build Public Website for production
```

### Database Operations

```bash
just rails-migrate    # Run Rails database migrations
just rails-seed       # Seed Rails database
just rails-setup      # Create, migrate, and seed Rails database
just rails-console    # Open Rails console
```

### Content API Operations

```bash
just content-api-build      # Build Content API binary
just content-api-watch      # Watch vault for changes (verbose mode)
just content-api-playground # Open GraphQL playground
just content-api-docker     # Run with Docker
```

## Architecture Overview

### Service Communication
- **website-ui** connects to both:
  - Content API at `http://localhost:1337/graphql` for page content (replacing Strapi)
  - Rails API at `http://localhost:3000/graphql` for project data
- **website-project-db-admin-ui2** connects to Rails API at `http://localhost:3000`
- Authentication handled via Google OAuth in the admin UI

### Key Technologies
- **Frontend**: Vue 3, Vue Router, Pinia (state management)
- **Backend**: Go with GraphQL (Content API), Rails 7 with GraphQL (Project API)
- **Databases**: SQLite for Rails, File-based storage for Content API
- **Content Management**: Obsidian vault with markdown files
- **Styling**: Tailwind CSS (website-ui), PrimeVue components
- **Testing**: Vitest (Vue apps), Rails built-in test framework, Go testing

### Data Models

#### Rails API Core Models
- **Project**: Main content entity with title, description, images
- **Category**: Hierarchical organization for tags
- **Tag**: Labeling system for projects
- **Event**: Recurring events associated with projects
- **Occurrence**: Individual instances of events
- **UmbrellaProject**: Grouping mechanism for related projects

#### Content API (Obsidian Markdown)
- **Pages**: Dynamic page content with component-based structure
- **Navigation**: Menu items and footer configuration
- **Assets**: Images and documents served statically

### GraphQL Schemas

#### Rails API (`/graphql`)
- Projects (with tags, events, occurrences)
- Categories (hierarchical structure)
- Search functionality
- Event management with recurrence patterns

#### Content API (`http://localhost:1337/graphql`)
- Pages with component-based content structure
- Navigation menus and footer
- Asset serving for images and documents

### Routing Patterns
- **website-ui**: Dynamic routes loaded from Content API pages
- **admin-ui**: Protected routes with Google OAuth authentication
- **Rails API**: RESTful + GraphQL endpoints under `/api/v1`
- **Content API**: GraphQL endpoint with static asset serving

## Important Configuration Files
- `website-project-db-api/config/database.yml` - Rails database config
- `website-ui/src/apollo/index.js` - Apollo client setup
- `website-project-db-admin-ui2/src/composables/auth_token.js` - Google OAuth authentication
- `website-content-api/obsidian-vault/` - Content files and structure
- `justfile` - Development task definitions

## Content Management

### Markdown Page Format (website-content-api)
Pages use a natural markdown format with component markers for easy editing in Obsidian:

```markdown
---
documentId: page-id
url: /page-url
title: Page Title
description: SEO description
author: Author Name
---

# [hero id="hero-main"] Hero Title

Description paragraph here.

[Learn More](/link)
![Image Alt](/uploads/image.jpg)

# [markdown id="content"]

## Regular markdown content

With any markdown formatting, headings, lists, etc.

# [highlightReel id="highlights"]

## Projects to Highlight
- project-1
- project-2

## Highlights

### Highlight Title
![Image Alt](/uploads/image.jpg)

# [contact id="contact"] Contact Name

Role or description
email@example.com

![Photo](/uploads/photo.jpg)

## Social Links
- [GitHub](https://github.com/username)
- [LinkedIn](https://linkedin.com/in/username)

# [projectSearch id="search" hide="internal,archive" tags="featured"]
```

### Key Features
- Components defined with `# [type id="..."]` syntax
- Natural markdown for content (no YAML strings)
- Type-specific parsers for complex components
- Backward compatible with YAML blocks
- Images: `![alt](url)`
- Links: `[text](url)`
- Internal notes can use HTML comments: `<!-- notes -->`

### Storage
- Pages stored in `website-content-api/obsidian-vault/pages/`
- Navigation in `website-content-api/obsidian-vault/navigation/`
- Hot-reload enabled for content changes during development
- Rails migrations in `website-project-db-api/db/migrate/`
- Active Storage configured for image uploads in Rails