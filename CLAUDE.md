# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monorepo containing four interconnected web applications for a Waldritter website system:

1. **website-pages** - Strapi v5 CMS for managing pages and content
2. **website-project-db-api** - Rails 7 API for project database with GraphQL support
3. **website-project-db-admin-ui2** - Vue 3 admin interface for managing projects
4. **website-ui** - Vue 3 public-facing website with Apollo GraphQL client

## Development Commands

### website-pages (Strapi CMS)
```bash
cd website-pages/website-pages
npm run develop    # Start dev server with auto-reload
npm run start      # Start production server
npm run build      # Build admin panel
npm run strapi     # Access Strapi CLI
```

### website-project-db-api (Rails API)
```bash
cd website-project-db-api
rails server       # Start development server on port 3000
rails console      # Open Rails console
rake test          # Run all tests
rake test:db       # Run tests with database reset
rake db:migrate    # Run database migrations
rake db:seed       # Seed database with initial data
```

### website-project-db-admin-ui2 (Admin UI)
```bash
cd website-project-db-admin-ui2
npm run dev        # Start Vite dev server
npm run build      # Build for production
npm run preview    # Preview production build
```

### website-ui (Public Website)
```bash
cd website-ui
pnpm dev           # Start Vite dev server
pnpm build         # Build for production
pnpm preview       # Preview production build
pnpm test:unit     # Run unit tests with Vitest
```

## Architecture Overview

### Service Communication
- **website-ui** connects to both:
  - Strapi CMS at `http://localhost:1337/graphql` for page content
  - Rails API at `http://localhost:3000/graphql` for project data
- **website-project-db-admin-ui2** connects to Rails API at `http://localhost:3000`
- Authentication handled via OIDC in the admin UI

### Key Technologies
- **Frontend**: Vue 3, Vue Router, Pinia (state management)
- **Backend**: Strapi v5 (Node.js), Rails 7 with GraphQL
- **Databases**: SQLite for both Strapi and Rails
- **Styling**: Tailwind CSS (website-ui), PrimeVue components
- **Testing**: Vitest (Vue apps), Rails built-in test framework

### Data Models

#### Rails API Core Models
- **Project**: Main content entity with title, description, images
- **Category**: Hierarchical organization for tags
- **Tag**: Labeling system for projects
- **Event**: Recurring events associated with projects
- **Occurrence**: Individual instances of events
- **UmbrellaProject**: Grouping mechanism for related projects

#### Strapi CMS Content Types
- **Page**: Dynamic page content with components
- **MenuItem**: Navigation structure
- **Footer**: Footer content management

### GraphQL Schema
The Rails API exposes a GraphQL endpoint at `/graphql` with types for:
- Projects (with tags, events, occurrences)
- Categories (hierarchical structure)
- Search functionality
- Event management with recurrence patterns

### Routing Patterns
- **website-ui**: Dynamic routes loaded from Strapi pages
- **admin-ui**: Protected routes with OIDC authentication
- **Rails API**: RESTful + GraphQL endpoints under `/api/v1`

## Important Configuration Files
- `website-project-db-api/config/database.yml` - Rails database config
- `website-ui/src/apollo/index.js` - Apollo client setup
- `website-project-db-admin-ui2/src/composables/auth_token.js` - OIDC authentication
- `website-pages/website-pages/config/` - Strapi configuration

## Database Operations
- Rails migrations in `website-project-db-api/db/migrate/`
- Strapi database managed automatically
- Active Storage configured for image uploads in Rails