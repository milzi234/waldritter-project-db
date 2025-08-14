package content

import (
	"strings"
)

// ProjectSearchParser parses projectSearch components from markdown
type ProjectSearchParser struct{}

// CanParse checks if this parser handles projectSearch components
func (p *ProjectSearchParser) CanParse(componentType string) bool {
	return componentType == "projectSearch"
}

// Parse converts markdown content into a ComponentProjectsProjectSearch
// Supports formats:
//
// Format 1 - Inline attributes:
// # [projectSearch id="filter" hide="internal,archive" tags="featured,current"]
//
// Format 2 - Section-based:
// # [projectSearch id="filter"]
// ## Filter Settings
// - Hide categories: internal, archive
// - Show tags: featured, current
func (p *ProjectSearchParser) Parse(heading string, content []string, attrs map[string]string) (Component, error) {
	component := &ComponentProjectsProjectSearch{
		ID: attrs["id"],
	}
	
	// Check for inline attributes first
	if hideCategories := attrs["hide"]; hideCategories != "" {
		if component.Filter == nil {
			component.Filter = &ProjectFilter{}
		}
		component.Filter.HiddenCategories = splitCommaSeparated(hideCategories)
	}
	
	if tags := attrs["tags"]; tags != "" {
		if component.Filter == nil {
			component.Filter = &ProjectFilter{}
		}
		component.Filter.SelectedTags = splitCommaSeparated(tags)
	}
	
	// Look for "Filter Settings" section
	filterSection := findSectionContent(content, "Filter Settings")
	if filterSection != nil {
		p.parseFilterSection(filterSection, component)
	} else {
		// Try parsing the entire content as filter settings
		p.parseFilterList(content, component)
	}
	
	return component, nil
}

// parseFilterSection parses filter settings from a section
func (p *ProjectSearchParser) parseFilterSection(lines []string, component *ComponentProjectsProjectSearch) {
	for _, line := range lines {
		if item, ok := parseListItem(line); ok {
			p.parseFilterItem(item, component)
		}
	}
}

// parseFilterList parses filter settings from content lines
func (p *ProjectSearchParser) parseFilterList(lines []string, component *ComponentProjectsProjectSearch) {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		
		// Check for list items
		if item, ok := parseListItem(line); ok {
			p.parseFilterItem(item, component)
		} else if strings.Contains(trimmed, ":") {
			// Also support non-list format like "Hide categories: x, y"
			p.parseFilterItem(trimmed, component)
		}
	}
}

// parseFilterItem parses a single filter configuration line
func (p *ProjectSearchParser) parseFilterItem(item string, component *ComponentProjectsProjectSearch) {
	lower := strings.ToLower(item)
	
	// Initialize filter if needed
	if component.Filter == nil {
		component.Filter = &ProjectFilter{}
	}
	
	// Check for hide categories
	if strings.Contains(lower, "hide") && strings.Contains(lower, "categor") {
		// Format: "Hide categories: internal, archive"
		parts := strings.SplitN(item, ":", 2)
		if len(parts) == 2 {
			categories := splitCommaSeparated(parts[1])
			component.Filter.HiddenCategories = append(component.Filter.HiddenCategories, categories...)
		}
	}
	
	// Check for show/selected tags
	if (strings.Contains(lower, "show") || strings.Contains(lower, "selected")) && strings.Contains(lower, "tag") {
		// Format: "Show tags: featured, current" or "Selected tags: featured"
		parts := strings.SplitN(item, ":", 2)
		if len(parts) == 2 {
			tags := splitCommaSeparated(parts[1])
			component.Filter.SelectedTags = append(component.Filter.SelectedTags, tags...)
		}
	}
}