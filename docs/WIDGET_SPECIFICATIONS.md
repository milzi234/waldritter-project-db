# Widget Specifications for WordPress Plugin

This document provides specifications for implementing WordPress widgets that integrate with the Waldritter Project Database API. These widgets replace the functionality previously provided by the Vue public website (`website-ui`) and Go content API (`website-content-api`).

## Data Sources

| Source | URL | Purpose |
|--------|-----|---------|
| **Project DB API** | `https://api.waldritter.de/graphql` | Projects, categories, tags, events, occurrences |
| **WordPress** | Native content | Static content (hero text, contact info, markdown) |

---

## Widget 1: Hero

**Purpose**: Large banner section with title, description, image, and call-to-action link.

**Data Source**: WordPress (static content)

### Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | Yes | Unique identifier for the component |
| `title` | string | Yes | Main heading text |
| `subtitle` | string | No | Secondary heading text |
| `description` | string | Yes | Body text/paragraph |
| `image` | Image | No | Background or hero image |
| `more_link` | string | No | URL for the call-to-action button |
| `link_text` | string | No | Text for the call-to-action button (default: "Learn More") |
| `isMarquee` | boolean | No | Display as scrolling marquee style |

### Image Object

```typescript
interface Image {
  url: string;       // Required - image URL
  caption?: string;  // Alt text/caption
  width?: number;    // Image width in pixels
  height?: number;   // Image height in pixels
}
```

### GraphQL Type (Reference)

```graphql
type ComponentPageHero {
  id: ID!
  title: String!
  description: String!
  more_link: String
  link_text: String
  image: Image
  isMarquee: Boolean
  subtitle: String
}
```

---

## Widget 2: Highlight Reel (Carousel)

**Purpose**: Carousel displaying featured projects and/or manual highlights.

**Data Source**: Project DB API (for projects) + WordPress (for manual highlights)

### Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `highlights` | Highlight[] | No | Manual highlight items (static content) |
| `projectsToHighlight` | string[] | No | Array of project IDs to feature |
| `projectFilter` | ProjectFilter | No | Filter to select projects dynamically |

### Highlight Object

```typescript
interface Highlight {
  title: string;       // Required - highlight title
  description: string; // Required - highlight description
  image: Image;        // Required - highlight image
  link?: string;       // URL to link to
  linkText?: string;   // Text for the link (default: "Homepage besuchen")
}
```

### ProjectFilter Object

```typescript
interface ProjectFilter {
  hiddenCategories?: string[];  // Category titles to exclude
  selectedTags?: string[];      // Tag titles to include (AND logic)
}
```

### Display Logic

1. If `projectsToHighlight` is provided, filter projects to only those IDs
2. If `projectFilter.selectedTags` is provided, include only projects with ALL specified tags
3. If `projectFilter.hiddenCategories` is provided, exclude projects in those categories
4. Combine manual `highlights` with filtered projects
5. If no filters specified and no manual highlights, show first 4 projects

### Carousel Configuration

```typescript
const responsiveOptions = [
  {
    breakpoint: '1024px',
    numVisible: 2,
    numScroll: 1
  },
  {
    breakpoint: '768px',
    numVisible: 1,
    numScroll: 1
  }
];
// Default: numVisible: 3, numScroll: 1
```

### GraphQL Type (Reference)

```graphql
type ComponentPageHighlightReel {
  id: ID!
  highlights: [Highlight!]!
  projectsToHighlight: [String!]
  projectFilter: ProjectFilter
}
```

---

## Widget 3: Project Search

**Purpose**: Interactive project browser with tag filtering.

**Data Source**: Project DB API

### Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `filter` | ProjectFilter | No | Initial filter configuration |

### Components

The Project Search widget contains two sub-components:

#### Tag Filter
- Displays categories and tags as checkboxes in a drawer/sidebar
- Multi-select checkboxes grouped by category
- Filter button shows count of selected tags: `Filter (3)`

#### Project List
- Displays filtered projects in a list/card layout
- Pagination with 5 items per page

### GraphQL Type (Reference)

```graphql
type ComponentProjectsProjectSearch {
  id: ID!
  filter: ProjectFilter
}

type ProjectFilter {
  hiddenCategories: [String!]
  selectedTags: [String!]
}
```

---

## Widget 4: Project List

**Purpose**: Display a list of projects with details.

**Data Source**: Project DB API

### Project Object

```typescript
interface Project {
  id: string;
  title: string;
  description: string;
  image: Image;
  homepage?: string;           // External website URL
  startDateTime?: string;      // ISO date string
  endDateTime?: string;        // ISO date string
  categories: Category[];
  occurrences: Occurrence[];   // Future event dates
}

interface Category {
  id: string;
  title: string;
  tags: Tag[];
}

interface Tag {
  id: string;
  title: string;
}

interface Occurrence {
  startDateTime: string;
  endDateTime?: string;
}
```

### Display Features

1. **Image**: Links to homepage if available, otherwise static
2. **Title**: Project name
3. **Time Range**: Start and end dates (if set)
4. **Description**: Project description text
5. **Tags**: Grouped by category (e.g., "Kategorie: Tag1, Tag2 | Status: Active")
6. **Homepage Link**: "Homepage besuchen" link if homepage URL exists
7. **Occurrences**: Expandable list of future event dates with pagination

### Pagination

- 5 projects per page (main list)
- 5 occurrences per page (nested list)

---

## Widget 5: Contact

**Purpose**: Display contact information with social links.

**Data Source**: WordPress (static content)

### Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Contact person/team name |
| `email` | string | Yes | Email address |
| `image` | Image | No | Contact photo |
| `social_links` | SocialLink[] | Yes | List of social media links |

### SocialLink Object

```typescript
interface SocialLink {
  label: string;     // Required - display text
  url: string;       // Required - link URL
  external: boolean; // Required - opens in new tab if true
  svgIcon?: string;  // Optional SVG icon markup
}
```

### GraphQL Type (Reference)

```graphql
type ComponentPageContact {
  id: ID!
  name: String!
  email: String!
  image: Image
  social_links: [SocialLink!]!
}

type SocialLink {
  label: String!
  url: String!
  external: Boolean!
  svgIcon: String
}
```

---

## Widget 6: Markdown

**Purpose**: Render arbitrary markdown content.

**Data Source**: WordPress (static content)

### Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `markdown` | string | Yes | Raw markdown content |

### Supported Markdown Features

- Headings (H1-H6)
- Paragraphs
- Bold, italic, strikethrough
- Lists (ordered and unordered)
- Links
- Images
- Code blocks
- Blockquotes
- Tables

### GraphQL Type (Reference)

```graphql
type ComponentPageMarkdown {
  id: ID!
  markdown: String!
}
```

---

## GraphQL Queries for Project DB API

### Fetch All Projects

```graphql
query GetProjects {
  projects {
    id
    title
    description
    homepage
    image {
      url
      caption
    }
    startDateTime
    endDateTime
    categories {
      id
      title
      tags {
        id
        title
      }
    }
    occurrences {
      startDateTime
      endDateTime
    }
  }
}
```

### Fetch Categories and Tags

```graphql
query GetCategories {
  categories {
    id
    title
    tags {
      id
      title
    }
  }
}
```

### Fetch Single Project

```graphql
query GetProject($id: ID!) {
  project(id: $id) {
    id
    title
    description
    homepage
    image {
      url
      caption
    }
    startDateTime
    endDateTime
    categories {
      id
      title
      tags {
        id
        title
      }
    }
    occurrences {
      startDateTime
      endDateTime
    }
  }
}
```

---

## Styling Guidelines

### Original CSS Framework
The original Vue implementation used:
- **Tailwind CSS** for utility classes
- **PrimeVue** components (Carousel, Paginator, Button, Drawer, Checkbox)

### Key Style Classes

```css
/* Card styling */
.bg-white.border.border-gray-300.rounded-lg.shadow-md.p-4.m-2

/* Image styling */
.w-full.h-48.object-cover.rounded-t-lg.mb-4

/* Title styling */
.text-xl.font-semibold.mb-2.text-gray-800

/* Link styling */
.text-green-600.hover:underline  /* Primary action links */
.text-green-600.hover:text-green-800  /* Secondary links */

/* Layout */
.flex.flex-col.md:flex-row  /* Responsive row/column */
.w-full.md:w-1/3  /* Responsive widths */
```

### Color Palette
- Primary: Green (`text-green-600`, `hover:text-green-800`)
- Text: Gray scale (`text-gray-600`, `text-gray-700`, `text-gray-800`)
- Background: White (`bg-white`)
- Borders: Light gray (`border-gray-300`)

---

## API Endpoints

| Service | URL | Authentication |
|---------|-----|----------------|
| Project DB API | `https://api.waldritter.de/graphql` | None (public read) |
| Admin API | `https://api.waldritter.de/graphql` | Google OAuth (write) |

---

## Notes for WordPress Implementation

1. **Caching**: Consider caching API responses as project data doesn't change frequently
2. **Error Handling**: Display graceful fallback if API is unavailable
3. **Responsive Design**: Ensure widgets work on mobile devices
4. **Accessibility**: Maintain proper ARIA labels and keyboard navigation
5. **SEO**: Ensure project content is indexable (server-side rendering recommended)
