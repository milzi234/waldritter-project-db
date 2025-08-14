# Component Parser Guide

## Overview

The website-content-api now supports type-specific markdown parsers that provide a more natural editing experience in Obsidian and other markdown editors. Components can be written using familiar markdown syntax instead of heavy YAML blocks.

## Component Types

### Hero Component

```markdown
# [hero id="main"] Welcome to Our Site

This is the hero description paragraph.

[Learn More](/about)
![Hero Banner](/uploads/hero.jpg)
```

**Parsed Elements:**
- Title: From the heading text after the component declaration
- Description: First paragraph(s) of content
- More Link: Markdown link with format `[text](url)`
- Image: Markdown image with format `![alt](url)`

### Highlight Reel Component

**Option 1: Section-based format**
```markdown
# [highlightReel id="featured"]

## Projects to Highlight
- project-slug-1
- project-slug-2

## Highlights

### First Feature Title
![Feature Image](/uploads/feature1.jpg)

### Second Feature Title
![Feature Image](/uploads/feature2.jpg)
```

**Option 2: Inline attributes with list format**
```markdown
# [highlightReel id="featured" projects="p1,p2,p3"]

- **Bold Title** - Description text
  ![Image](/uploads/image.jpg)

- **Another Title** - Another description
  ![Image](/uploads/image2.jpg)
```

**Parsed Elements:**
- Projects: From list items under "Projects to Highlight" section or `projects` attribute
- Highlights: H3 headings with their associated images, or list items with bold text
- Filter settings: Can use `hide="category1,category2"` and `tags="tag1,tag2"` attributes

### Project Search Component

```markdown
# [projectSearch id="filter" hide="internal,archive" tags="featured,current"]
```

Or with content:

```markdown
# [projectSearch id="filter"]

## Filter Settings
- Hide categories: internal, archive
- Show tags: featured, current
```

**Parsed Elements:**
- Hidden categories: From `hide` attribute or "Hide categories:" lines
- Selected tags: From `tags` attribute or "Show tags:" lines

### Contact Component

```markdown
# [contact id="team-lead"] John Doe

Lead Organizer
john.doe@waldritter.de

![John Doe](/uploads/john.jpg)

## Social Links
- [GitHub](https://github.com/johndoe)
- [LinkedIn](https://linkedin.com/in/johndoe)
```

**Parsed Elements:**
- Name: From heading text or first line of content
- Email: Line containing @ symbol
- Image: Markdown image
- Social links: Markdown links, automatically detects platform icons

### Regular Markdown Sections

Any H1 heading without a component type creates a markdown component:

```markdown
# Regular Section Title

This content becomes a regular markdown component.
Any markdown syntax works here.

## Subsection

More content...
```

## YAML Override

YAML blocks still work and take precedence over markdown parsing:

```markdown
# [hero id="test"] Markdown Title

This markdown description will be ignored.

` ` `yaml hero
title: "YAML Title (takes precedence)"
description: "YAML Description (takes precedence)"
` ` `
```

## Parser Features

### Intelligent Section Detection

The parser intelligently detects when sections belong to a component:
- H2 and H3 headings within a component are treated as sections
- H1 headings always start a new component
- Content is associated with the nearest component or section

### Attribute Support

Component declarations support inline attributes:
```markdown
# [componentType id="my-id" attr1="value1" attr2="value2"] Optional Title
```

### Backward Compatibility

- All existing YAML-based content continues to work
- YAML blocks override markdown parsing when present
- You can mix markdown and YAML formats in the same file

## Benefits

1. **Natural Editing**: Write content as you would normally write markdown
2. **Better Preview**: Obsidian and other editors can preview the content meaningfully
3. **Less Boilerplate**: No need for complex YAML structures for simple components
4. **Flexibility**: Use markdown for simple cases, YAML for complex data
5. **Maintainability**: Easier to read and edit content files

## Implementation Details

### Parser Architecture

```
ComponentParser (interface)
├── HeroParser
├── HighlightReelParser
├── ProjectSearchParser
├── ContactParser
└── MarkdownParser
```

Each parser:
- Implements `CanParse(componentType string) bool`
- Implements `Parse(heading, content, attrs) (Component, error)`
- Handles specific markdown patterns for its component type

### Helper Functions

- `parseMarkdownImage()`: Extracts images from `![alt](url)` syntax
- `parseMarkdownLink()`: Extracts links from `[text](url)` syntax
- `extractBoldText()`: Extracts text from `**bold**` markers
- `findSectionContent()`: Finds content under specific headings
- `splitCommaSeparated()`: Parses comma-separated attribute values

## Testing

The parser system includes comprehensive tests:
- Unit tests for each parser
- Integration tests for mixed content
- Backward compatibility tests
- Edge case handling

Run tests with:
```bash
go test ./content -v
```

## Migration Guide

### Converting YAML to Markdown Format

**Before (YAML):**
```yaml
# [hero id="main"]

` ` `yaml hero
title: "Welcome"
description: "Our description"
more_link: "/about"
image:
  url: "/hero.jpg"
  caption: "Hero Image"
` ` `
```

**After (Markdown):**
```markdown
# [hero id="main"] Welcome

Our description

[Learn More](/about)
![Hero Image](/hero.jpg)
```

The new format is more concise and natural to write while maintaining all functionality.