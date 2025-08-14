package content

import (
	"strings"
)

// ContactParser parses contact components from markdown
type ContactParser struct{}

// CanParse checks if this parser handles contact components
func (p *ContactParser) CanParse(componentType string) bool {
	return componentType == "contact"
}

// Parse converts markdown content into a ComponentPageContact
// Supports format:
// # [contact id="team-lead"] John Doe
// Lead Organizer
// john.doe@waldritter.de
// ![John Doe](/uploads/john.jpg)
// ## Social Links
// - [GitHub](https://github.com/johndoe)
// - [LinkedIn](https://linkedin.com/in/johndoe)
func (p *ContactParser) Parse(heading string, content []string, attrs map[string]string) (Component, error) {
	component := &ComponentPageContact{
		ID:          attrs["id"],
		Name:        strings.TrimSpace(heading),
		SocialLinks: make([]*SocialLink, 0),
	}
	
	// Look for Social Links section first
	socialSection := findSectionContent(content, "Social Links")
	if socialSection != nil {
		p.parseSocialLinks(socialSection, component)
	}
	
	// Parse main content (before any sections)
	mainContent := p.extractMainContent(content)
	p.parseMainContent(mainContent, component)
	
	return component, nil
}

// extractMainContent gets content before any H2 sections
func (p *ContactParser) extractMainContent(lines []string) []string {
	var mainContent []string
	for _, line := range lines {
		level, _ := isHeading(line)
		if level == 2 {
			// Stop at first H2
			break
		}
		mainContent = append(mainContent, line)
	}
	return mainContent
}

// parseMainContent extracts name, role, email, and image from main content
func (p *ContactParser) parseMainContent(lines []string, component *ComponentPageContact) {
	foundEmail := false
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		
		// Check for image
		if img := parseMarkdownImage(trimmed); img != nil {
			component.Image = img
			continue
		}
		
		// Check for email (contains @)
		if !foundEmail && strings.Contains(trimmed, "@") {
			component.Email = trimmed
			foundEmail = true
			continue
		}
		
		// If name is empty and this is the first content line, use it as name
		if component.Name == "" && i == 0 {
			component.Name = trimmed
		}
	}
}

// parseSocialLinks extracts social links from a section
func (p *ContactParser) parseSocialLinks(lines []string, component *ComponentPageContact) {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Check for list item with link
		if item, ok := parseListItem(line); ok {
			// Check if the item contains a markdown link
			if text, url := parseMarkdownLink(item); text != "" && url != "" {
				socialLink := &SocialLink{
					Label:    text,
					URL:      url,
					External: p.isExternalURL(url),
				}
				
				// Try to detect icon type from label or URL
				socialLink.SVGIcon = p.detectSocialIcon(text, url)
				
				component.SocialLinks = append(component.SocialLinks, socialLink)
			}
		} else if text, url := parseMarkdownLink(trimmed); text != "" && url != "" {
			// Handle links not in list format
			socialLink := &SocialLink{
				Label:    text,
				URL:      url,
				External: p.isExternalURL(url),
			}
			socialLink.SVGIcon = p.detectSocialIcon(text, url)
			component.SocialLinks = append(component.SocialLinks, socialLink)
		}
	}
}

// isExternalURL checks if a URL is external
func (p *ContactParser) isExternalURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// detectSocialIcon tries to detect the appropriate icon based on label or URL
func (p *ContactParser) detectSocialIcon(label, url string) string {
	lower := strings.ToLower(label + " " + url)
	
	// Common social platforms
	if strings.Contains(lower, "github") {
		return "github"
	}
	if strings.Contains(lower, "linkedin") {
		return "linkedin"
	}
	if strings.Contains(lower, "twitter") || strings.Contains(lower, "x.com") {
		return "twitter"
	}
	if strings.Contains(lower, "facebook") {
		return "facebook"
	}
	if strings.Contains(lower, "instagram") {
		return "instagram"
	}
	if strings.Contains(lower, "youtube") {
		return "youtube"
	}
	if strings.Contains(lower, "email") || strings.Contains(lower, "mailto") {
		return "email"
	}
	
	return ""
}