package content

import (
	"regexp"
	"strings"
)

// ComponentParser defines the interface for type-specific component parsers
type ComponentParser interface {
	// CanParse checks if this parser can handle the given component type
	CanParse(componentType string) bool
	
	// Parse converts markdown content into a Component
	// heading is the original heading text (after the component type)
	// content is the content lines under the heading
	// attrs are any attributes extracted from the heading
	Parse(heading string, content []string, attrs map[string]string) (Component, error)
}

// ComponentParserRegistry manages all available component parsers
type ComponentParserRegistry struct {
	parsers []ComponentParser
}

// NewComponentParserRegistry creates a new registry with all default parsers
func NewComponentParserRegistry() *ComponentParserRegistry {
	return &ComponentParserRegistry{
		parsers: []ComponentParser{
			&HeroParser{},
			&HighlightReelParser{},
			&ProjectSearchParser{},
			&ContactParser{},
			&MarkdownParser{},
		},
	}
}

// GetParser returns the appropriate parser for a component type
func (r *ComponentParserRegistry) GetParser(componentType string) ComponentParser {
	for _, parser := range r.parsers {
		if parser.CanParse(componentType) {
			return parser
		}
	}
	return nil
}

// Helper functions for parsing markdown elements

// parseMarkdownImage extracts image information from markdown syntax
// Supports: ![alt text](url)
func parseMarkdownImage(line string) *Image {
	re := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 3 {
		return &Image{
			URL:     matches[2],
			Caption: matches[1],
		}
	}
	return nil
}

// parseMarkdownLink extracts link text and URL from markdown syntax
// Supports: [text](url)
func parseMarkdownLink(line string) (text, url string) {
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 3 {
		return matches[1], matches[2]
	}
	return "", ""
}

// extractBoldText extracts text between ** markers
func extractBoldText(line string) string {
	re := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// parseListItem parses a markdown list item (- or *)
func parseListItem(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "- ") {
		return strings.TrimSpace(trimmed[2:]), true
	}
	if strings.HasPrefix(trimmed, "* ") {
		return strings.TrimSpace(trimmed[2:]), true
	}
	return "", false
}

// splitCommaSeparated splits a comma-separated string and trims each item
func splitCommaSeparated(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// isHeading checks if a line is a markdown heading and returns its level
func isHeading(line string) (level int, text string) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "######") {
		return 6, strings.TrimSpace(trimmed[6:])
	}
	if strings.HasPrefix(trimmed, "#####") {
		return 5, strings.TrimSpace(trimmed[5:])
	}
	if strings.HasPrefix(trimmed, "####") {
		return 4, strings.TrimSpace(trimmed[4:])
	}
	if strings.HasPrefix(trimmed, "###") {
		return 3, strings.TrimSpace(trimmed[3:])
	}
	if strings.HasPrefix(trimmed, "##") {
		return 2, strings.TrimSpace(trimmed[2:])
	}
	if strings.HasPrefix(trimmed, "#") {
		return 1, strings.TrimSpace(trimmed[1:])
	}
	return 0, ""
}

// findSectionContent extracts content under a specific heading
func findSectionContent(lines []string, sectionName string) []string {
	var content []string
	inSection := false
	sectionLevel := 0
	foundSection := false
	
	for _, line := range lines {
		level, text := isHeading(line)
		if level > 0 {
			if inSection && level <= sectionLevel {
				// Found next section at same or higher level, stop
				break
			}
			if strings.EqualFold(text, sectionName) {
				inSection = true
				sectionLevel = level
				foundSection = true
				continue
			}
		}
		
		if inSection {
			content = append(content, line)
		}
	}
	
	// Return nil if section was not found, empty slice if section exists but has no content
	if !foundSection {
		return nil
	}
	return content
}

// extractAttributeValue extracts a value from attributes or returns empty string
func extractAttributeValue(attrs map[string]string, keys ...string) string {
	for _, key := range keys {
		if val, ok := attrs[key]; ok {
			return val
		}
	}
	return ""
}