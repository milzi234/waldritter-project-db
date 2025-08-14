package content

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

// Parser handles markdown and YAML parsing
type Parser struct {
	markdown        goldmark.Markdown
	componentParsers *ComponentParserRegistry
}

// NewParser creates a new parser instance
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	return &Parser{
		markdown:         md,
		componentParsers: NewComponentParserRegistry(),
	}
}

// ExtractFrontmatter separates YAML frontmatter from markdown content
func ExtractFrontmatter(content []byte) (frontmatter []byte, markdown []byte) {
	contentStr := string(content)
	
	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(contentStr, "---") {
		return []byte{}, content
	}
	
	// Find the closing delimiter
	// Look for "\n---\n" or "\n---\r\n"
	idx := strings.Index(contentStr[3:], "\n---")
	if idx == -1 {
		// No closing delimiter found
		return []byte{}, content
	}
	
	// Adjust index to account for the offset
	idx += 3
	
	// Extract frontmatter (between the delimiters)
	fm := contentStr[3:idx]
	fm = strings.TrimSpace(fm)
	
	// Find where content starts after the closing ---
	contentStart := idx + 4 // Skip "\n---"
	if contentStart < len(contentStr) && contentStr[contentStart] == '\n' {
		contentStart++ // Skip the newline after ---
	}
	
	// Extract content
	var contentPart string
	if contentStart < len(contentStr) {
		contentPart = strings.TrimSpace(contentStr[contentStart:])
	}
	
	return []byte(fm), []byte(contentPart)
}

// ParsePageMetadata parses YAML frontmatter into PageMetadataYAML
func ParsePageMetadata(frontmatter []byte) (*PageMetadataYAML, error) {
	var meta PageMetadataYAML
	
	if len(frontmatter) == 0 {
		return &meta, nil
	}
	
	err := yaml.Unmarshal(frontmatter, &meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page metadata: %w", err)
	}
	
	return &meta, nil
}

// componentContext holds the current component being parsed
type componentContext struct {
	componentType string
	id            string
	title         string
	attributes    map[string]string
	content       []string
	yamlData      map[string]interface{}
}

// ParseMarkdownPage parses markdown content to extract components
func (p *Parser) ParseMarkdownPage(content []byte) ([]Component, error) {
	if len(content) == 0 {
		return []Component{}, nil
	}
	
	lines := strings.Split(string(content), "\n")
	components := []Component{}
	var currentComponent *componentContext
	inYamlBlock := false
	yamlBlockType := ""
	yamlLines := []string{}
	
	// Regular expressions for parsing
	headingWithTypeRe := regexp.MustCompile(`^(#{1,6})\s*\[(\w+)(?:\s+([^\]]+))?\]\s*(.*)$`)
	attributeRe := regexp.MustCompile(`(\w+)="([^"]*)"`)
	yamlBlockStartRe := regexp.MustCompile("^```yaml\\s+(\\w+)\\s*$")
	yamlBlockEndRe := regexp.MustCompile("^```\\s*$")
	
	for i, line := range lines {
		// Check for YAML block start
		if matches := yamlBlockStartRe.FindStringSubmatch(line); matches != nil {
			inYamlBlock = true
			yamlBlockType = matches[1]
			yamlLines = []string{}
			continue
		}
		
		// Check for YAML block end
		if inYamlBlock && yamlBlockEndRe.MatchString(line) {
			inYamlBlock = false
			
			// Parse YAML and merge with current component
			if currentComponent != nil && yamlBlockType == currentComponent.componentType {
				yamlContent := strings.Join(yamlLines, "\n")
				var yamlData map[string]interface{}
				if err := yaml.Unmarshal([]byte(yamlContent), &yamlData); err == nil {
					// Merge YAML data with current component
					for k, v := range yamlData {
						currentComponent.yamlData[k] = v
					}
				}
			}
			yamlBlockType = ""
			yamlLines = []string{}
			continue
		}
		
		// Collect YAML block content
		if inYamlBlock {
			yamlLines = append(yamlLines, line)
			continue
		}
		
		// Check for heading with component type
		if matches := headingWithTypeRe.FindStringSubmatch(line); matches != nil {
			// Save previous component if exists
			if currentComponent != nil {
				if comp := p.buildComponent(currentComponent); comp != nil {
					components = append(components, comp)
				}
			}
			
			// Start new component
			componentType := matches[2]
			attributesStr := matches[3]
			title := strings.TrimSpace(matches[4])
			
			// Parse attributes
			attributes := make(map[string]string)
			if attributesStr != "" {
				for _, match := range attributeRe.FindAllStringSubmatch(attributesStr, -1) {
					attributes[match[1]] = match[2]
				}
			}
			
			currentComponent = &componentContext{
				componentType: componentType,
				id:            attributes["id"],
				title:         title,
				attributes:    attributes,
				content:       []string{},
				yamlData:      make(map[string]interface{}),
			}
			continue
		}
		
		// Check for regular H1 heading (which should start a new markdown component)
		if strings.HasPrefix(line, "# ") && !headingWithTypeRe.MatchString(line) {
			// H1 heading without component type - always starts a new markdown component
			// Save previous component if exists
			if currentComponent != nil {
				if comp := p.buildComponent(currentComponent); comp != nil {
					components = append(components, comp)
				}
			}
			
			currentComponent = &componentContext{
				componentType: "markdown",
				id:            "markdown-" + strconv.Itoa(i),
				content:       []string{line},
				yamlData:      make(map[string]interface{}),
				attributes:    make(map[string]string),
			}
		} else if currentComponent != nil {
			// Add content to current component (including lower-level headings)
			currentComponent.content = append(currentComponent.content, line)
		} else if strings.HasPrefix(line, "#") && !headingWithTypeRe.MatchString(line) {
			// Non-H1 heading without current component - create a markdown component
			currentComponent = &componentContext{
				componentType: "markdown",
				id:            "markdown-" + strconv.Itoa(i),
				content:       []string{line},
				yamlData:      make(map[string]interface{}),
				attributes:    make(map[string]string),
			}
		}
	}
	
	// Save last component
	if currentComponent != nil {
		if comp := p.buildComponent(currentComponent); comp != nil {
			components = append(components, comp)
		}
	}
	
	return components, nil
}

// buildComponent creates a Component from componentContext
func (p *Parser) buildComponent(ctx *componentContext) Component {
	if ctx == nil {
		return nil
	}
	
	// Clean up content - trim empty lines from start and end
	content := ctx.content
	for len(content) > 0 && strings.TrimSpace(content[0]) == "" {
		content = content[1:]
	}
	for len(content) > 0 && strings.TrimSpace(content[len(content)-1]) == "" {
		content = content[:len(content)-1]
	}
	
	// If we have YAML data, it takes precedence - use the existing YAML-based parsing
	if len(ctx.yamlData) > 0 {
		return p.buildComponentFromYAML(ctx, content)
	}
	
	// Try to use a type-specific parser if available
	if parser := p.componentParsers.GetParser(ctx.componentType); parser != nil {
		// Make sure ID is in attributes if it's not already there
		attrs := ctx.attributes
		if ctx.id != "" && attrs["id"] == "" {
			attrs = make(map[string]string)
			for k, v := range ctx.attributes {
				attrs[k] = v
			}
			attrs["id"] = ctx.id
		}
		if comp, err := parser.Parse(ctx.title, content, attrs); err == nil {
			return comp
		}
		// Fall back to YAML-based parsing if type-specific parsing fails
	}
	
	// Fall back to YAML-based parsing for unknown types
	return p.buildComponentFromYAML(ctx, content)
}

// buildComponentFromYAML creates a Component using the original YAML-based approach
func (p *Parser) buildComponentFromYAML(ctx *componentContext, content []string) Component {
	if ctx == nil {
		return nil
	}
	
	switch ctx.componentType {
	case "hero":
		component := &ComponentPageHero{
			ID:    ctx.id,
			Title: ctx.title,
		}
		
		// Extract description from first paragraph
		if len(content) > 0 {
			// Find first non-empty line as description
			for _, line := range content {
				if strings.TrimSpace(line) != "" {
					component.Description = strings.TrimSpace(line)
					break
				}
			}
		}
		
		// Override with YAML data if present
		if desc, ok := ctx.yamlData["description"].(string); ok {
			component.Description = desc
		}
		if title, ok := ctx.yamlData["title"].(string); ok {
			component.Title = title
		}
		if moreLink, ok := ctx.yamlData["more_link"].(string); ok {
			component.MoreLink = moreLink
		}
		
		// Handle image from YAML
		if imageData, ok := ctx.yamlData["image"].(map[string]interface{}); ok {
			component.Image = &Image{}
			if url, ok := imageData["url"].(string); ok {
				component.Image.URL = url
			}
			if caption, ok := imageData["caption"].(string); ok {
				component.Image.Caption = caption
			}
			if width, ok := imageData["width"].(int); ok {
				component.Image.Width = width
			}
			if height, ok := imageData["height"].(int); ok {
				component.Image.Height = height
			}
		}
		
		return component
		
	case "markdown":
		// Join all content lines for markdown
		markdownContent := strings.Join(content, "\n")
		
		// Override with YAML data if present
		if md, ok := ctx.yamlData["markdown"].(string); ok {
			markdownContent = md
		} else if cont, ok := ctx.yamlData["content"].(string); ok {
			markdownContent = cont
		}
		
		return &ComponentPageMarkdown{
			ID:       ctx.id,
			Markdown: markdownContent,
		}
		
	case "contact":
		component := &ComponentPageContact{
			ID:          ctx.id,
			SocialLinks: make([]*SocialLink, 0),
		}
		
		// Extract name and email from YAML
		if name, ok := ctx.yamlData["name"].(string); ok {
			component.Name = name
		}
		if email, ok := ctx.yamlData["email"].(string); ok {
			component.Email = email
		}
		
		// Handle image from YAML
		if imageData, ok := ctx.yamlData["image"].(map[string]interface{}); ok {
			component.Image = &Image{}
			if url, ok := imageData["url"].(string); ok {
				component.Image.URL = url
			}
			if caption, ok := imageData["caption"].(string); ok {
				component.Image.Caption = caption
			}
		}
		
		// Handle social links from YAML
		if socialLinks, ok := ctx.yamlData["social_links"].([]interface{}); ok {
			for _, linkData := range socialLinks {
				if link, ok := linkData.(map[string]interface{}); ok {
					socialLink := &SocialLink{}
					if label, ok := link["label"].(string); ok {
						socialLink.Label = label
					}
					if url, ok := link["url"].(string); ok {
						socialLink.URL = url
					}
					if external, ok := link["external"].(bool); ok {
						socialLink.External = external
					}
					if svgIcon, ok := link["svgIcon"].(string); ok {
						socialLink.SVGIcon = svgIcon
					}
					component.SocialLinks = append(component.SocialLinks, socialLink)
				}
			}
		}
		
		return component
		
	case "projectSearch":
		component := &ComponentProjectsProjectSearch{
			ID: ctx.id,
		}
		
		// Handle filter from YAML
		if filterData, ok := ctx.yamlData["filter"].(map[string]interface{}); ok {
			component.Filter = &ProjectFilter{}
			if hidden, ok := filterData["hiddenCategories"].([]interface{}); ok {
				for _, cat := range hidden {
					if catStr, ok := cat.(string); ok {
						component.Filter.HiddenCategories = append(component.Filter.HiddenCategories, catStr)
					}
				}
			}
			if selected, ok := filterData["selectedTags"].([]interface{}); ok {
				for _, tag := range selected {
					if tagStr, ok := tag.(string); ok {
						component.Filter.SelectedTags = append(component.Filter.SelectedTags, tagStr)
					}
				}
			}
		}
		
		return component
		
	case "highlightReel":
		component := &ComponentPageHighlightReel{
			ID:         ctx.id,
			Highlights: make([]*Highlight, 0),
		}
		
		// Handle highlights from YAML
		if highlights, ok := ctx.yamlData["highlights"].([]interface{}); ok {
			for _, hlData := range highlights {
				if hl, ok := hlData.(map[string]interface{}); ok {
					highlight := &Highlight{}
					if desc, ok := hl["description"].(string); ok {
						highlight.Description = desc
					}
					if imageData, ok := hl["image"].(map[string]interface{}); ok {
						highlight.Image = &Image{}
						if url, ok := imageData["url"].(string); ok {
							highlight.Image.URL = url
						}
						if caption, ok := imageData["caption"].(string); ok {
							highlight.Image.Caption = caption
						}
					}
					component.Highlights = append(component.Highlights, highlight)
				}
			}
		}
		
		// Handle projectsToHighlight from YAML
		if projects, ok := ctx.yamlData["projectsToHighlight"].([]interface{}); ok {
			for _, proj := range projects {
				if projStr, ok := proj.(string); ok {
					component.ProjectsToHighlight = append(component.ProjectsToHighlight, projStr)
				}
			}
		}
		
		// Handle projectFilter from YAML
		if filterData, ok := ctx.yamlData["projectFilter"].(map[string]interface{}); ok {
			component.ProjectFilter = &ProjectFilter{}
			if hidden, ok := filterData["hiddenCategories"].([]interface{}); ok {
				for _, cat := range hidden {
					if catStr, ok := cat.(string); ok {
						component.ProjectFilter.HiddenCategories = append(component.ProjectFilter.HiddenCategories, catStr)
					}
				}
			}
			if selected, ok := filterData["selectedTags"].([]interface{}); ok {
				for _, tag := range selected {
					if tagStr, ok := tag.(string); ok {
						component.ProjectFilter.SelectedTags = append(component.ProjectFilter.SelectedTags, tagStr)
					}
				}
			}
		}
		
		return component
	}
	
	return nil
}

// ParseMenuMetadata parses YAML frontmatter into MenuMetadata
func ParseMenuMetadata(frontmatter []byte) (*MenuMetadata, error) {
	var meta MenuMetadata
	
	err := yaml.Unmarshal(frontmatter, &meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse menu metadata: %w", err)
	}
	
	return &meta, nil
}

// ParseFooterMetadata parses YAML frontmatter into FooterMetadata
func ParseFooterMetadata(frontmatter []byte) (*FooterMetadata, error) {
	var meta FooterMetadata
	
	err := yaml.Unmarshal(frontmatter, &meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse footer metadata: %w", err)
	}
	
	return &meta, nil
}

// ParseComponent converts a ComponentDefinition into a GraphQL component type
func (p *Parser) ParseComponent(def ComponentDefinition) (Component, error) {
	switch def.Type {
	case "markdown":
		// Use Content field if Markdown field is empty
		markdownContent := def.Markdown
		if markdownContent == "" {
			markdownContent = def.Content
		}
		return &ComponentPageMarkdown{
			ID:       def.ID,
			Markdown: markdownContent,
		}, nil
		
	case "hero":
		component := &ComponentPageHero{
			ID:          def.ID,
			Title:       def.Title,
			Description: def.Description,
			MoreLink:    def.MoreLink,
		}
		
		if def.Image != nil {
			component.Image = &Image{
				URL:     def.Image.URL,
				Caption: def.Image.Caption,
				Width:   def.Image.Width,
				Height:  def.Image.Height,
			}
		}
		
		return component, nil
		
	case "contact":
		component := &ComponentPageContact{
			ID:          def.ID,
			Name:        def.Name,
			Email:       def.Email,
			SocialLinks: make([]*SocialLink, 0),
		}
		
		if def.Image != nil {
			component.Image = &Image{
				URL:     def.Image.URL,
				Caption: def.Image.Caption,
				Width:   def.Image.Width,
				Height:  def.Image.Height,
			}
		}
		
		for _, link := range def.SocialLinks {
			component.SocialLinks = append(component.SocialLinks, &SocialLink{
				Label:    link.Label,
				URL:      link.URL,
				External: link.External,
				SVGIcon:  link.SVGIcon,
			})
		}
		
		return component, nil
		
	case "projectSearch":
		component := &ComponentProjectsProjectSearch{
			ID: def.ID,
		}
		
		if def.Filter != nil {
			component.Filter = &ProjectFilter{
				HiddenCategories: def.Filter.HiddenCategories,
				SelectedTags:     def.Filter.SelectedTags,
			}
		}
		
		return component, nil
		
	case "highlightReel":
		component := &ComponentPageHighlightReel{
			ID:                  def.ID,
			Highlights:          make([]*Highlight, 0),
			ProjectsToHighlight: def.ProjectsToHighlight,
		}
		
		for _, hl := range def.Highlights {
			highlight := &Highlight{
				Description: hl.Description,
			}
			
			if hl.Image != nil {
				highlight.Image = &Image{
					URL:     hl.Image.URL,
					Caption: hl.Image.Caption,
					Width:   hl.Image.Width,
					Height:  hl.Image.Height,
				}
			}
			
			component.Highlights = append(component.Highlights, highlight)
		}
		
		if def.ProjectFilter != nil {
			component.ProjectFilter = &ProjectFilter{
				HiddenCategories: def.ProjectFilter.HiddenCategories,
				SelectedTags:     def.ProjectFilter.SelectedTags,
			}
		}
		
		return component, nil
		
	default:
		return nil, fmt.Errorf("unknown component type: %s", def.Type)
	}
}

// RenderMarkdown converts markdown to HTML
func (p *Parser) RenderMarkdown(markdown string) (string, error) {
	if markdown == "" {
		return "", nil
	}
	
	var buf bytes.Buffer
	if err := p.markdown.Convert([]byte(markdown), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}
	
	return buf.String(), nil
}

// ConvertMenuItems converts MenuItemDefinition slice to MenuItem slice for GraphQL
func ConvertMenuItems(items []MenuItemDefinition) []*MenuItem {
	result := make([]*MenuItem, 0, len(items))
	
	for _, item := range items {
		menuItem := &MenuItem{
			DocumentID: item.DocumentID,
			Label:      item.Label,
			To:         item.To,
			Order:      item.Order,
		}
		
		if len(item.Items) > 0 {
			menuItem.Items = ConvertMenuItems(item.Items)
		}
		
		result = append(result, menuItem)
	}
	
	return result
}

// ConvertFooter converts FooterMetadata to Footer for GraphQL
func ConvertFooter(meta *FooterMetadata) *Footer {
	footer := &Footer{
		ShortDescriptionTitle: meta.ShortDescriptionTitle,
		ShortDescriptionText:  meta.ShortDescriptionText,
		Links1:                make([]*FooterLink, 0, len(meta.Links1)),
		Links2:                make([]*FooterLink, 0, len(meta.Links2)),
		SocialLinks:           make([]*FooterSocialLink, 0, len(meta.SocialLinks)),
	}
	
	for _, link := range meta.Links1 {
		footer.Links1 = append(footer.Links1, &FooterLink{
			Label: link.Label,
			URL:   link.URL,
		})
	}
	
	for _, link := range meta.Links2 {
		footer.Links2 = append(footer.Links2, &FooterLink{
			Label: link.Label,
			URL:   link.URL,
		})
	}
	
	for _, link := range meta.SocialLinks {
		footer.SocialLinks = append(footer.SocialLinks, &FooterSocialLink{
			Label:   link.Label,
			URL:     link.URL,
			SVGIcon: link.SVGIcon,
		})
	}
	
	return footer
}

// ParseMarkdownMenu parses markdown headings to extract menu structure
func ParseMarkdownMenu(markdown []byte) ([]*MenuItem, error) {
	// Convert to string and split into lines
	content := string(markdown)
	lines := strings.Split(content, "\n")
	
	// Track the menu hierarchy
	var rootItems []*MenuItem
	var currentH1 *MenuItem
	var currentH2 *MenuItem
	var h1Order, h2Order, h3Order int
	hasMainRoot := false
	
	// Process each line
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check for H1 heading
		if strings.HasPrefix(line, "# ") {
			label, url := extractLabelAndURL(strings.TrimPrefix(line, "# "))
			
			// Check if this is the $MAIN root
			if label == "$MAIN" {
				hasMainRoot = true
				// Skip creating a menu item for $MAIN, but reset state for its children
				currentH1 = nil
				currentH2 = nil
				h1Order = 0  // Reset order counter for direct children of $MAIN
				h2Order = 0
				h3Order = 0
				continue
			}
			
			// If we're not under $MAIN or $MAIN hasn't been declared, treat as normal H1
			if !hasMainRoot {
				// Create H1 menu item
				h1Order++
				item := &MenuItem{
					DocumentID: generateDocumentID(label, ""),
					Label:      label,
					To:         url,
					Order:      h1Order,
					Items:      []*MenuItem{},
				}
				
				rootItems = append(rootItems, item)
				currentH1 = item
				currentH2 = nil
				h2Order = 0
				h3Order = 0
			}
			
		} else if strings.HasPrefix(line, "## ") {
			// H2 heading
			label, url := extractLabelAndURL(strings.TrimPrefix(line, "## "))
			
			if hasMainRoot {
				// If we're under $MAIN, all H2s are top-level items
				h1Order++
				item := &MenuItem{
					DocumentID: generateDocumentID(label, ""),
					Label:      label,
					To:         url,
					Order:      h1Order,
					Items:      []*MenuItem{},
				}
				rootItems = append(rootItems, item)
				currentH1 = item
				currentH2 = nil
				h2Order = 0
				h3Order = 0
			} else if currentH1 != nil {
				// Normal H2 under H1 (when not under $MAIN)
				h2Order++
				h3Order = 0
				parentPrefix := getParentPrefix(currentH1.DocumentID)
				item := &MenuItem{
					DocumentID: generateDocumentID(label, parentPrefix),
					Label:      label,
					To:         url,
					Order:      h2Order,
					Items:      []*MenuItem{},
				}
				currentH1.Items = append(currentH1.Items, item)
				currentH2 = item
			}
			
		} else if strings.HasPrefix(line, "### ") && currentH2 != nil {
			// H3 heading - only process if we have a parent H2
			label, url := extractLabelAndURL(strings.TrimPrefix(line, "### "))
			
			h3Order++
			parentPrefix := getParentPrefix(currentH2.DocumentID)
			item := &MenuItem{
				DocumentID: generateDocumentID(label, parentPrefix),
				Label:      label,
				To:         url,
				Order:      h3Order,
				Items:      []*MenuItem{},
			}
			
			currentH2.Items = append(currentH2.Items, item)
		}
		// Ignore all other lines (non-heading content)
	}
	
	// Wrap in $MAIN root if not already present
	mainRoot := &MenuItem{
		DocumentID: "main-root",
		Label:      "$MAIN",
		Items:      rootItems,
	}
	
	return []*MenuItem{mainRoot}, nil
}

// extractLabelAndURL extracts the label and optional URL from a heading line
func extractLabelAndURL(text string) (label string, url string) {
	text = strings.TrimSpace(text)
	
	// Check for URL in brackets at the end
	if idx := strings.LastIndex(text, "["); idx > 0 {
		if endIdx := strings.LastIndex(text, "]"); endIdx > idx {
			label = strings.TrimSpace(text[:idx])
			url = strings.TrimSpace(text[idx+1 : endIdx])
			return label, url
		}
	}
	
	// No URL found, entire text is the label
	return text, ""
}

// generateDocumentID creates a document ID from a label
func generateDocumentID(label string, parentPrefix string) string {
	// Convert label to lowercase and replace special characters
	id := strings.ToLower(label)
	
	// Replace common special characters
	replacements := map[string]string{
		" ":  "-",
		"ä":  "ae",
		"ö":  "oe",
		"ü":  "ue",
		"ß":  "ss",
		"é":  "e",
		"è":  "e",
		"ê":  "e",
		"ë":  "e",
		"à":  "a",
		"â":  "a",
		"ô":  "o",
		"î":  "i",
		"ï":  "i",
		"ç":  "c",
		"&":  "-",
		"@":  "-",
		"!":  "",
		"%":  "",
		"'":  "",
		"\"": "",
		".":  "",
		",":  "",
		":":  "",
		";":  "",
		"?":  "",
		"/":  "-",
		"\\": "-",
		"(":  "",
		")":  "",
		"[":  "",
		"]":  "",
		"{":  "",
		"}":  "",
		"<":  "",
		">":  "",
		"|":  "-",
		"+":  "-",
		"=":  "-",
		"*":  "",
		"#":  "",
		"$":  "",
		"^":  "",
		"~":  "",
		"`":  "",
	}
	
	for old, new := range replacements {
		id = strings.ReplaceAll(id, old, new)
	}
	
	// Remove multiple consecutive hyphens
	for strings.Contains(id, "--") {
		id = strings.ReplaceAll(id, "--", "-")
	}
	
	// Trim hyphens from start and end
	id = strings.Trim(id, "-")
	
	// Add menu prefix
	if parentPrefix != "" {
		id = parentPrefix + "-" + id
	} else {
		id = "menu-" + id
	}
	
	return id
}

// getParentPrefix extracts the parent prefix from a document ID
func getParentPrefix(documentID string) string {
	// For nested items, use the parent's full ID as prefix
	return documentID
}