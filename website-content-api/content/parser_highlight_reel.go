package content

import (
	"strings"
)

// HighlightReelParser parses highlightReel components from markdown
type HighlightReelParser struct{}

// CanParse checks if this parser handles highlightReel components
func (p *HighlightReelParser) CanParse(componentType string) bool {
	return componentType == "highlightReel"
}

// Parse converts markdown content into a ComponentPageHighlightReel
// Supports multiple formats:
//
// Format 1 - Section-based:
// # [highlightReel id="featured"]
// ## Projects to Highlight
// - project1
// - project2
// ## Highlights
// ### Title of highlight
// ![Description](/image.jpg)
//
// Format 2 - Inline attributes with list:
// # [highlightReel id="featured" projects="project1,project2"]
// - **Title** - Description
//   ![Image](/image.jpg)
func (p *HighlightReelParser) Parse(heading string, content []string, attrs map[string]string) (Component, error) {
	component := &ComponentPageHighlightReel{
		ID:         attrs["id"],
		Highlights: make([]*Highlight, 0),
	}
	
	// Check for inline projects attribute
	if projects := attrs["projects"]; projects != "" {
		component.ProjectsToHighlight = splitCommaSeparated(projects)
	}
	
	// Check for inline filter attributes
	if hideCategories := attrs["hide"]; hideCategories != "" {
		if component.ProjectFilter == nil {
			component.ProjectFilter = &ProjectFilter{}
		}
		component.ProjectFilter.HiddenCategories = splitCommaSeparated(hideCategories)
	}
	if tags := attrs["tags"]; tags != "" {
		if component.ProjectFilter == nil {
			component.ProjectFilter = &ProjectFilter{}
		}
		component.ProjectFilter.SelectedTags = splitCommaSeparated(tags)
	}
	
	// Look for "Projects to Highlight" section
	projectsSection := findSectionContent(content, "Projects to Highlight")
	if projectsSection != nil && len(component.ProjectsToHighlight) == 0 {
		for _, line := range projectsSection {
			if item, ok := parseListItem(line); ok {
				component.ProjectsToHighlight = append(component.ProjectsToHighlight, item)
			}
		}
	}
	
	// Look for "Highlights" section
	highlightsSection := findSectionContent(content, "Highlights")
	if highlightsSection != nil {
		p.parseHighlightsSection(highlightsSection, component)
	} else {
		// Try parsing the entire content as a list of highlights
		p.parseHighlightsList(content, component)
	}
	
	return component, nil
}

// parseHighlightsSection parses highlights from a section with H3 headings
func (p *HighlightReelParser) parseHighlightsSection(lines []string, component *ComponentPageHighlightReel) {
	var currentHighlight *Highlight
	
	for _, line := range lines {
		level, text := isHeading(line)
		
		// H3 starts a new highlight
		if level == 3 {
			if currentHighlight != nil {
				component.Highlights = append(component.Highlights, currentHighlight)
			}
			currentHighlight = &Highlight{
				Description: text,
			}
			continue
		}
		
		// Check for image
		if currentHighlight != nil {
			if img := parseMarkdownImage(line); img != nil {
				currentHighlight.Image = img
			}
		}
	}
	
	// Don't forget the last highlight
	if currentHighlight != nil {
		component.Highlights = append(component.Highlights, currentHighlight)
	}
}

// parseHighlightsList parses highlights from a markdown list
func (p *HighlightReelParser) parseHighlightsList(lines []string, component *ComponentPageHighlightReel) {
	var currentHighlight *Highlight
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Check if this is a list item
		if item, ok := parseListItem(line); ok {
			// Save previous highlight if exists
			if currentHighlight != nil {
				component.Highlights = append(component.Highlights, currentHighlight)
			}
			
			// Create new highlight
			currentHighlight = &Highlight{}
			
			// Check for bold text as title with description
			if boldText := extractBoldText(item); boldText != "" {
				// Format: **Title** - Description
				parts := strings.SplitN(item, " - ", 2)
				if len(parts) == 2 {
					currentHighlight.Description = strings.TrimSpace(parts[1])
				} else {
					// Just use the whole text
					currentHighlight.Description = item
				}
			} else {
				currentHighlight.Description = item
			}
			
			// Check if next line is an indented image
			if i+1 < len(lines) {
				nextLine := lines[i+1]
				// Check for indented content (2+ spaces or tab)
				if strings.HasPrefix(nextLine, "  ") || strings.HasPrefix(nextLine, "\t") {
					if img := parseMarkdownImage(strings.TrimSpace(nextLine)); img != nil {
						currentHighlight.Image = img
					}
				}
			}
		} else if currentHighlight != nil && strings.HasPrefix(line, "  ") {
			// Indented content under a list item
			if img := parseMarkdownImage(trimmed); img != nil {
				currentHighlight.Image = img
			}
		}
	}
	
	// Don't forget the last highlight
	if currentHighlight != nil {
		component.Highlights = append(component.Highlights, currentHighlight)
	}
}