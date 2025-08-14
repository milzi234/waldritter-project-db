# Markdown Page Format Documentation

The website-content-api now supports a more natural Obsidian-friendly markdown format for defining page components.

## New Format Structure

Pages now consist of:
1. **YAML Frontmatter** - Contains metadata (documentId, url, title, description, author, etc.)
2. **Markdown Body** - Contains component definitions using special heading syntax

## Component Definition Syntax

### Basic Component

```markdown
# [componentType id="component-id"] Component Title

Component content goes here...
```

### Component with YAML Configuration

````markdown
# [componentType id="component-id"] Component Title

Component content...

```yaml componentType
property: value
nested:
  property: value
```
````

## Supported Component Types

### Hero Component

```markdown
# [hero id="main-hero"] Welcome to Our Site

This paragraph becomes the hero description.

```yaml hero
more_link: /learn-more
image:
  url: /assets/hero-image.jpg
  caption: Hero Image Caption
```
```

### Markdown Component

```markdown
# [markdown id="content-section"] Content Section

## Subheading

All content after the component heading becomes markdown content.

- List items
- **Bold text**
- *Italic text*

Regular markdown syntax is preserved.
```

### Contact Component

```markdown
# [contact id="contact-form"] Get in Touch

```yaml contact
name: John Doe
email: john@example.com
image:
  url: /assets/john.jpg
social_links:
  - label: Twitter
    url: https://twitter.com/johndoe
    external: true
  - label: LinkedIn
    url: https://linkedin.com/in/johndoe
    external: true
```
```

### Project Search Component

```markdown
# [projectSearch id="search"] Find Projects

```yaml projectSearch
filter:
  hiddenCategories:
    - internal
    - deprecated
  selectedTags:
    - featured
    - public
```
```

### Highlight Reel Component

```markdown
# [highlightReel id="highlights"] Featured Work

```yaml highlightReel
highlights:
  - description: First project highlight
    image:
      url: /assets/project1.jpg
      caption: Project 1
  - description: Second project highlight
    image:
      url: /assets/project2.jpg
      caption: Project 2
projectsToHighlight:
  - project-id-1
  - project-id-2
projectFilter:
  selectedTags:
    - featured
```
```

## Complete Example Page

```markdown
---
documentId: home
url: /
title: Home Page
description: Welcome to our website
author: Site Admin
---

# [hero id="hero-main"] Welcome to Our Company

We build amazing digital experiences that delight users and drive results.

```yaml hero
more_link: /about
image:
  url: /assets/hero-home.jpg
  caption: Our office space
```

# [markdown id="intro"] Introduction

## What We Do

We specialize in creating modern web applications using cutting-edge technologies.

Our expertise includes:
- Full-stack development
- Cloud architecture
- DevOps and CI/CD
- User experience design

## Our Approach

We believe in **iterative development** and *continuous improvement*.

# [highlightReel id="featured-work"] Our Latest Projects

```yaml highlightReel
projectsToHighlight:
  - project-alpha
  - project-beta
  - project-gamma
projectFilter:
  selectedTags:
    - featured
    - recent
```

# [contact id="contact-cta"] Let's Work Together

```yaml contact
name: Sales Team
email: sales@example.com
social_links:
  - label: Twitter
    url: https://twitter.com/ourcompany
    external: true
  - label: GitHub
    url: https://github.com/ourcompany
    external: true
```
```

## Migration from Old Format

### Old Format (YAML Frontmatter)
```yaml
---
documentId: page-id
url: /page-url
components:
  - type: hero
    id: hero-1
    title: Welcome
    description: Description
  - type: markdown
    id: content-1
    content: |
      # Heading
      Content here
---
```

### New Format (Markdown Body)
```markdown
---
documentId: page-id
url: /page-url
---

# [hero id="hero-1"] Welcome

Description

# [markdown id="content-1"] Heading

Content here
```

## Key Benefits

1. **Obsidian-Friendly**: Content is visible and editable in Obsidian's editor
2. **Natural Writing**: Write content naturally without escaping in YAML
3. **Preview Support**: Component content previews correctly in Obsidian
4. **Flexible Structure**: Mix regular markdown with component definitions
5. **YAML Override**: Optional YAML blocks for complex configuration

## Parsing Rules

1. **Component Boundaries**: Only headings with `[type]` syntax create new components
2. **Regular Headings**: Headings without `[type]` become part of current component or start a new markdown component
3. **Content Collection**: All content between component headings belongs to that component
4. **YAML Merging**: YAML blocks merge with the current component if types match
5. **Attribute Parsing**: Attributes in heading (e.g., `id="value"`) are parsed and applied

## Notes

- Component IDs should be unique within a page
- YAML blocks must specify the component type they belong to
- Empty lines between components are trimmed
- The parser preserves markdown formatting within components
- Regular markdown content without component markers creates a markdown component automatically