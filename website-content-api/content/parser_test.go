package content

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		wantFrontmatter string
		wantContent     string
		wantError       bool
	}{
		{
			name: "valid frontmatter with content",
			input: []byte(`---
documentId: test-page
url: /test
title: Test Page
---

# Content

This is the content`),
			wantFrontmatter: `documentId: test-page
url: /test
title: Test Page`,
			wantContent: `# Content

This is the content`,
			wantError: false,
		},
		{
			name: "no frontmatter",
			input: []byte(`# Just Content

No frontmatter here`),
			wantFrontmatter: "",
			wantContent: `# Just Content

No frontmatter here`,
			wantError: false,
		},
		{
			name: "empty frontmatter",
			input: []byte(`---
---

# Content`),
			wantFrontmatter: "",
			wantContent: `# Content`,
			wantError: false,
		},
		{
			name: "frontmatter only",
			input: []byte(`---
documentId: test
url: /test
---`),
			wantFrontmatter: `documentId: test
url: /test`,
			wantContent: "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frontmatter, content := ExtractFrontmatter(tt.input)
			
			if string(frontmatter) != tt.wantFrontmatter {
				t.Errorf("ExtractFrontmatter() frontmatter = %q, want %q", string(frontmatter), tt.wantFrontmatter)
			}
			
			if string(content) != tt.wantContent {
				t.Errorf("ExtractFrontmatter() content = %q, want %q", string(content), tt.wantContent)
			}
		})
	}
}

func TestParsePageMetadata(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *PageMetadataYAML
		wantErr bool
	}{
		{
			name: "complete metadata",
			input: []byte(`
documentId: about-page
url: /about
title: "About Us"
description: "Learn about our team"
author: "John Doe"
publishDate: 2024-01-15`),
			want: &PageMetadataYAML{
				DocumentID:  "about-page",
				URL:         "/about",
				Title:       "About Us",
				Description: "Learn about our team",
				Author:      "John Doe",
			},
			wantErr: false,
		},
		{
			name: "minimal metadata",
			input: []byte(`
documentId: test
url: /test`),
			want: &PageMetadataYAML{
				DocumentID: "test",
				URL:        "/test",
			},
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			input:   []byte(`invalid: [yaml`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePageMetadata(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePageMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if got.DocumentID != tt.want.DocumentID {
					t.Errorf("ParsePageMetadata() DocumentID = %v, want %v", got.DocumentID, tt.want.DocumentID)
				}
				if got.URL != tt.want.URL {
					t.Errorf("ParsePageMetadata() URL = %v, want %v", got.URL, tt.want.URL)
				}
				if got.Title != tt.want.Title {
					t.Errorf("ParsePageMetadata() Title = %v, want %v", got.Title, tt.want.Title)
				}
			}
		})
	}
}

func TestParseComponent(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		input   ComponentDefinition
		want    interface{}
		wantErr bool
	}{
		{
			name: "markdown component",
			input: ComponentDefinition{
				Type:    "markdown",
				ID:      "md-1",
				Content: "# Hello\n\nThis is **markdown**",
			},
			want: &ComponentPageMarkdown{
				ID:       "md-1",
				Markdown: "# Hello\n\nThis is **markdown**",
			},
			wantErr: false,
		},
		{
			name: "hero component",
			input: ComponentDefinition{
				Type:        "hero",
				ID:          "hero-1",
				Title:       "Welcome",
				Description: "Hero description",
				MoreLink:    "/learn-more",
				Image: &ImageDefinition{
					URL:     "/assets/hero.jpg",
					Caption: "Hero image",
				},
			},
			want: &ComponentPageHero{
				ID:          "hero-1",
				Title:       "Welcome",
				Description: "Hero description",
				MoreLink:    "/learn-more",
				Image: &Image{
					URL:     "/assets/hero.jpg",
					Caption: "Hero image",
				},
			},
			wantErr: false,
		},
		{
			name: "contact component",
			input: ComponentDefinition{
				Type:  "contact",
				ID:    "contact-1",
				Name:  "John Doe",
				Email: "john@example.com",
				SocialLinks: []SocialLinkDefinition{
					{
						Label:    "Twitter",
						URL:      "https://twitter.com/johndoe",
						External: true,
					},
				},
			},
			want: &ComponentPageContact{
				ID:    "contact-1",
				Name:  "John Doe",
				Email: "john@example.com",
				SocialLinks: []*SocialLink{
					{
						Label:    "Twitter",
						URL:      "https://twitter.com/johndoe",
						External: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "project search component",
			input: ComponentDefinition{
				Type: "projectSearch",
				ID:   "search-1",
				Filter: &FilterDefinition{
					HiddenCategories: []string{"internal"},
					SelectedTags:     []string{"featured"},
				},
			},
			want: &ComponentProjectsProjectSearch{
				ID: "search-1",
				Filter: &ProjectFilter{
					HiddenCategories: []string{"internal"},
					SelectedTags:     []string{"featured"},
				},
			},
			wantErr: false,
		},
		{
			name: "unknown component type",
			input: ComponentDefinition{
				Type: "unknown",
				ID:   "unknown-1",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseComponent(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseComponent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got == nil {
				t.Errorf("ParseComponent() returned nil component")
			}
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "simple markdown",
			input:   "# Hello\n\nThis is **bold** and *italic*",
			want:    "<h1>Hello</h1>\n<p>This is <strong>bold</strong> and <em>italic</em></p>\n",
			wantErr: false,
		},
		{
			name:    "markdown with links",
			input:   "[Link](https://example.com)",
			want:    "<p><a href=\"https://example.com\">Link</a></p>\n",
			wantErr: false,
		},
		{
			name:    "markdown with code",
			input:   "```go\nfunc main() {}\n```",
			want:    "<pre><code class=\"language-go\">func main() {}\n</code></pre>\n",
			wantErr: false,
		},
		{
			name:    "empty markdown",
			input:   "",
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.RenderMarkdown(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if got != tt.want {
				t.Errorf("RenderMarkdown() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseMenuMetadata(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *MenuMetadata
		wantErr bool
	}{
		{
			name: "valid menu structure",
			input: []byte(`
type: menu
documentId: main-navigation
menuItems:
  - documentId: main-root
    label: "$MAIN"
    items:
      - documentId: menu-home
        label: "Home"
        to: "/"
        order: 1
      - documentId: menu-about
        label: "About"
        to: "/about"
        order: 2
        items:
          - documentId: menu-team
            label: "Team"
            to: "/about/team"
            order: 1`),
			want: &MenuMetadata{
				Type:       "menu",
				DocumentID: "main-navigation",
				MenuItems: []MenuItemDefinition{
					{
						DocumentID: "main-root",
						Label:      "$MAIN",
						Items: []MenuItemDefinition{
							{
								DocumentID: "menu-home",
								Label:      "Home",
								To:         "/",
								Order:      1,
							},
							{
								DocumentID: "menu-about",
								Label:      "About",
								To:         "/about",
								Order:      2,
								Items: []MenuItemDefinition{
									{
										DocumentID: "menu-team",
										Label:      "Team",
										To:         "/about/team",
										Order:      1,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			input:   []byte(`invalid: [yaml`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMenuMetadata(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMenuMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				if got.DocumentID != tt.want.DocumentID {
					t.Errorf("ParseMenuMetadata() DocumentID = %v, want %v", got.DocumentID, tt.want.DocumentID)
				}
				if len(got.MenuItems) != len(tt.want.MenuItems) {
					t.Errorf("ParseMenuMetadata() MenuItems length = %v, want %v", len(got.MenuItems), len(tt.want.MenuItems))
				}
			}
		})
	}
}

func TestParseFooterMetadata(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *FooterMetadata
		wantErr bool
	}{
		{
			name: "valid footer structure",
			input: []byte(`
type: footer
documentId: main-footer
shortDescriptionTitle: "Test Site"
shortDescriptionText: "This is a test site"
links_1:
  - label: "About"
    url: "/about"
  - label: "Contact"
    url: "/contact"
links_2:
  - label: "Privacy"
    url: "/privacy"
social_links:
  - label: "Twitter"
    url: "https://twitter.com/test"
    svgIcon: "<svg>...</svg>"`),
			want: &FooterMetadata{
				Type:                  "footer",
				DocumentID:            "main-footer",
				ShortDescriptionTitle: "Test Site",
				ShortDescriptionText:  "This is a test site",
				Links1: []FooterLinkDefinition{
					{Label: "About", URL: "/about"},
					{Label: "Contact", URL: "/contact"},
				},
				Links2: []FooterLinkDefinition{
					{Label: "Privacy", URL: "/privacy"},
				},
				SocialLinks: []FooterSocialLinkDefinition{
					{Label: "Twitter", URL: "https://twitter.com/test", SVGIcon: "<svg>...</svg>"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFooterMetadata(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFooterMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && got != nil {
				if got.ShortDescriptionTitle != tt.want.ShortDescriptionTitle {
					t.Errorf("ParseFooterMetadata() ShortDescriptionTitle = %v, want %v", got.ShortDescriptionTitle, tt.want.ShortDescriptionTitle)
				}
				if len(got.Links1) != len(tt.want.Links1) {
					t.Errorf("ParseFooterMetadata() Links1 length = %v, want %v", len(got.Links1), len(tt.want.Links1))
				}
			}
		})
	}
}

func TestParseMarkdownMenu(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []*MenuItem
		wantErr bool
	}{
		{
			name: "simple flat menu structure",
			input: []byte(`# Start [/]

# Leitbild [/leitbild]

# Termine [/termine]

# Kontakt [/kontakt]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-start",
							Label:      "Start",
							To:         "/",
							Order:      1,
						},
						{
							DocumentID: "menu-leitbild",
							Label:      "Leitbild",
							To:         "/leitbild",
							Order:      2,
						},
						{
							DocumentID: "menu-termine",
							Label:      "Termine",
							To:         "/termine",
							Order:      3,
						},
						{
							DocumentID: "menu-kontakt",
							Label:      "Kontakt",
							To:         "/kontakt",
							Order:      4,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nested menu items (2 levels)",
			input: []byte(`# Start [/]

# Schulangebote

## Grundschule [/schulangebote/grundschule]

## Sekundarstufe I [/schulangebote/sekundarstufe-1]

## Sekundarstufe II [/schulangebote/sekundarstufe-2]

# Kontakt [/kontakt]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-start",
							Label:      "Start",
							To:         "/",
							Order:      1,
						},
						{
							DocumentID: "menu-schulangebote",
							Label:      "Schulangebote",
							Order:      2,
							Items: []*MenuItem{
								{
									DocumentID: "menu-schulangebote-grundschule",
									Label:      "Grundschule",
									To:         "/schulangebote/grundschule",
									Order:      1,
								},
								{
									DocumentID: "menu-schulangebote-sekundarstufe-i",
									Label:      "Sekundarstufe I",
									To:         "/schulangebote/sekundarstufe-1",
									Order:      2,
								},
								{
									DocumentID: "menu-schulangebote-sekundarstufe-ii",
									Label:      "Sekundarstufe II",
									To:         "/schulangebote/sekundarstufe-2",
									Order:      3,
								},
							},
						},
						{
							DocumentID: "menu-kontakt",
							Label:      "Kontakt",
							To:         "/kontakt",
							Order:      3,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "deeply nested menu items (3 levels)",
			input: []byte(`# Products

## Software

### Development Tools [/products/software/dev-tools]

### Office Suite [/products/software/office]

## Hardware [/products/hardware]

# About [/about]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-products",
							Label:      "Products",
							Order:      1,
							Items: []*MenuItem{
								{
									DocumentID: "menu-products-software",
									Label:      "Software",
									Order:      1,
									Items: []*MenuItem{
										{
											DocumentID: "menu-products-software-development-tools",
											Label:      "Development Tools",
											To:         "/products/software/dev-tools",
											Order:      1,
										},
										{
											DocumentID: "menu-products-software-office-suite",
											Label:      "Office Suite",
											To:         "/products/software/office",
											Order:      2,
										},
									},
								},
								{
									DocumentID: "menu-products-hardware",
									Label:      "Hardware",
									To:         "/products/hardware",
									Order:      2,
								},
							},
						},
						{
							DocumentID: "menu-about",
							Label:      "About",
							To:         "/about",
							Order:      2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "menu with special characters in labels",
			input: []byte(`# Über Uns [/ueber-uns]

# Café & Restaurant [/cafe-restaurant]

# 50% Off! [/sale]

# Contact@Email [/contact]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-ueber-uns",
							Label:      "Über Uns",
							To:         "/ueber-uns",
							Order:      1,
						},
						{
							DocumentID: "menu-cafe-restaurant",
							Label:      "Café & Restaurant",
							To:         "/cafe-restaurant",
							Order:      2,
						},
						{
							DocumentID: "menu-50-off",
							Label:      "50% Off!",
							To:         "/sale",
							Order:      3,
						},
						{
							DocumentID: "menu-contact-email",
							Label:      "Contact@Email",
							To:         "/contact",
							Order:      4,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty content",
			input: []byte(``),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items:      []*MenuItem{},
				},
			},
			wantErr: false,
		},
		{
			name: "whitespace only content",
			input: []byte(`

  

`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items:      []*MenuItem{},
				},
			},
			wantErr: false,
		},
		{
			name: "menu with text content between headings",
			input: []byte(`# Home [/]

This is some text that should be ignored.

# About [/about]

More text here that should also be ignored.

## Team [/about/team]

## Mission [/about/mission]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-home",
							Label:      "Home",
							To:         "/",
							Order:      1,
						},
						{
							DocumentID: "menu-about",
							Label:      "About",
							To:         "/about",
							Order:      2,
							Items: []*MenuItem{
								{
									DocumentID: "menu-about-team",
									Label:      "Team",
									To:         "/about/team",
									Order:      1,
								},
								{
									DocumentID: "menu-about-mission",
									Label:      "Mission",
									To:         "/about/mission",
									Order:      2,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "menu with mixed link and no-link items",
			input: []byte(`# Home [/]

# Products

## Hardware [/products/hardware]

## Software

### Desktop Apps [/products/software/desktop]

### Web Apps [/products/software/web]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-home",
							Label:      "Home",
							To:         "/",
							Order:      1,
						},
						{
							DocumentID: "menu-products",
							Label:      "Products",
							Order:      2,
							Items: []*MenuItem{
								{
									DocumentID: "menu-products-hardware",
									Label:      "Hardware",
									To:         "/products/hardware",
									Order:      1,
								},
								{
									DocumentID: "menu-products-software",
									Label:      "Software",
									Order:      2,
									Items: []*MenuItem{
										{
											DocumentID: "menu-products-software-desktop-apps",
											Label:      "Desktop Apps",
											To:         "/products/software/desktop",
											Order:      1,
										},
										{
											DocumentID: "menu-products-software-web-apps",
											Label:      "Web Apps",
											To:         "/products/software/web",
											Order:      2,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "menu with invalid heading jump (H1 to H3)",
			input: []byte(`# Home [/]

### This should be treated as regular text, not a menu item

# About [/about]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-home",
							Label:      "Home",
							To:         "/",
							Order:      1,
						},
						{
							DocumentID: "menu-about",
							Label:      "About",
							To:         "/about",
							Order:      2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "menu with $MAIN root already present",
			input: []byte(`# $MAIN

## Home [/]

## About [/about]`),
			want: []*MenuItem{
				{
					DocumentID: "main-root",
					Label:      "$MAIN",
					Items: []*MenuItem{
						{
							DocumentID: "menu-home",
							Label:      "Home",
							To:         "/",
							Order:      1,
						},
						{
							DocumentID: "menu-about",
							Label:      "About",
							To:         "/about",
							Order:      2,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMarkdownMenu(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMarkdownMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if !compareMenuItems(got, tt.want) {
					t.Errorf("ParseMarkdownMenu() mismatch\ngot:  %v\nwant: %v", 
						formatMenuItems(got), formatMenuItems(tt.want))
				}
			}
		})
	}
}

// Helper function to compare menu items recursively
func compareMenuItems(a, b []*MenuItem) bool {
	if len(a) != len(b) {
		return false
	}
	
	for i := range a {
		if !compareMenuItem(a[i], b[i]) {
			return false
		}
	}
	
	return true
}

func compareMenuItem(a, b *MenuItem) bool {
	if a == nil || b == nil {
		return a == b
	}
	
	if a.DocumentID != b.DocumentID ||
		a.Label != b.Label ||
		a.To != b.To ||
		a.Order != b.Order {
		return false
	}
	
	return compareMenuItems(a.Items, b.Items)
}

// Helper function to format menu items for error messages
func formatMenuItems(items []*MenuItem) string {
	return formatMenuItemsIndent(items, 0)
}

func formatMenuItemsIndent(items []*MenuItem, indent int) string {
	result := ""
	prefix := strings.Repeat("  ", indent)
	
	for _, item := range items {
		result += fmt.Sprintf("%s[%s] %s", prefix, item.DocumentID, item.Label)
		if item.To != "" {
			result += fmt.Sprintf(" -> %s", item.To)
		}
		result += fmt.Sprintf(" (order: %d)\n", item.Order)
		
		if len(item.Items) > 0 {
			result += formatMenuItemsIndent(item.Items, indent+1)
		}
	}
	
	return result
}