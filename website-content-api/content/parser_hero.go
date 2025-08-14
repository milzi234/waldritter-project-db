package content

import (
	"strings"
)

// HeroParser parses hero components from markdown
type HeroParser struct{}

// CanParse checks if this parser handles hero components
func (p *HeroParser) CanParse(componentType string) bool {
	return componentType == "hero"
}

// Parse converts markdown content into a ComponentPageHero
// Supports formats:
// # [hero id="main"] Title Here
// Description paragraph here.
// [Learn More](/about)
// ![Hero Image](/uploads/hero.jpg)
func (p *HeroParser) Parse(heading string, content []string, attrs map[string]string) (Component, error) {
	component := &ComponentPageHero{
		ID:    attrs["id"],
		Title: strings.TrimSpace(heading),
	}
	
	// Parse content for description, link, and image
	for i, line := range content {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		
		// Check for image
		if img := parseMarkdownImage(trimmed); img != nil {
			component.Image = img
			continue
		}
		
		// Check for link (Learn More)
		if text, url := parseMarkdownLink(trimmed); text != "" && url != "" {
			component.MoreLink = url
			continue
		}
		
		// First non-empty, non-special line is the description
		if component.Description == "" && !strings.HasPrefix(trimmed, "#") {
			// For multi-line descriptions, collect until we hit a blank line or special element
			descLines := []string{trimmed}
			for j := i + 1; j < len(content); j++ {
				nextLine := strings.TrimSpace(content[j])
				if nextLine == "" {
					break
				}
				// Stop if we hit an image or link
				if strings.HasPrefix(nextLine, "![") || strings.HasPrefix(nextLine, "[") {
					break
				}
				descLines = append(descLines, nextLine)
			}
			component.Description = strings.Join(descLines, " ")
		}
	}
	
	return component, nil
}