package content

import (
	"strings"
)

// MarkdownParser parses markdown components
type MarkdownParser struct{}

// CanParse checks if this parser handles markdown components
func (p *MarkdownParser) CanParse(componentType string) bool {
	return componentType == "markdown"
}

// Parse converts markdown content into a ComponentPageMarkdown
func (p *MarkdownParser) Parse(heading string, content []string, attrs map[string]string) (Component, error) {
	// For markdown components, just join all content lines
	markdownContent := strings.Join(content, "\n")
	
	// Trim empty lines from start and end
	markdownContent = strings.TrimSpace(markdownContent)
	
	component := &ComponentPageMarkdown{
		ID:       attrs["id"],
		Markdown: markdownContent,
	}
	
	// If there's a heading and no content, use the heading as content
	if markdownContent == "" && heading != "" {
		component.Markdown = heading
	}
	
	return component, nil
}