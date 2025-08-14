---
documentId: example-page
url: /example
title: Example Page with Natural Markdown Parsing
description: Demonstrates the new type-specific parsers for a more markdown-native editing experience
---

# [hero id="main-hero"] Welcome to Waldritter

Join us in creating epic LARP adventures that bring communities together through storytelling and roleplay.

[Learn More](/about)
![Epic LARP Adventure](/uploads/hero-banner.jpg)

# [highlightReel id="featured-highlights"]

## Projects to Highlight
- larp-fuer-demokratie
- sommerfreizeit-2024
- neue-ortsgruppe-berlin

## Highlights

### 🎭 LARP für Demokratie
![Democracy through Roleplay](/uploads/democracy-larp.jpg)

### 🏕️ Sommerfreizeit 2024
![Summer Camp Adventures](/uploads/summer-camp.jpg)

### 🌟 10 Ortsgruppen bundesweit
![Local Groups Map](/uploads/groups-map.jpg)

# [projectSearch id="project-filter" hide="internal,archive" tags="featured,current"]

# [contact id="contact-person"] Max Mustermann

Projektleiter Waldritter e.V.
max.mustermann@waldritter.de

![Max Mustermann](/uploads/max.jpg)

## Social Links
- [GitHub](https://github.com/waldritter)
- [Instagram](https://instagram.com/waldritter)
- [Facebook](https://facebook.com/waldritter)

# Alternative Formats

You can also use more compact inline formats:

# [highlightReel id="compact" projects="project1,project2,project3"]

- **Quick Setup** - Get started in minutes
  ![Quick Setup](/uploads/setup.jpg)

- **Community Driven** - Built by LARPers for LARPers
  ![Community](/uploads/community.jpg)

- **Open Source** - Free and open for everyone
  ![Open Source](/uploads/opensource.jpg)

# Regular Markdown Content

This section is just regular markdown that will be parsed as a markdown component.

## Features

- Type-specific parsers for natural markdown editing
- Backward compatible with YAML blocks
- Intelligent section detection
- Support for nested content

## Benefits

1. **More Natural Editing**: Write in markdown, not YAML
2. **Better Obsidian Experience**: Leverage Obsidian's markdown preview
3. **Flexible**: Mix markdown and YAML as needed

# YAML Override Example

Sometimes you still want to use YAML for complex data:

# [hero id="yaml-hero"] This title will be overridden

```yaml hero
title: "YAML Override Title"
description: "When YAML is present, it takes precedence over markdown parsing"
more_link: "/yaml-example"
image:
  url: "/uploads/yaml-image.jpg"
  caption: "YAML-defined image"
  width: 1920
  height: 1080
```

This demonstrates that YAML blocks still work and override markdown parsing when present.