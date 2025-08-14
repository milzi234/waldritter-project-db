package content

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseMarkdownPage(t *testing.T) {
	parser := NewParser()
	
	tests := []struct {
		name    string
		input   string
		want    []Component
		wantErr bool
	}{
		{
			name: "hero component with title and description",
			input: `# [hero id="hero-main"] Welcome to Our Site

This is the hero description that appears below the title.

More content here.`,
			want: []Component{
				&ComponentPageHero{
					ID:          "hero-main",
					Title:       "Welcome to Our Site",
					Description: "This is the hero description that appears below the title.",
				},
			},
			wantErr: false,
		},
		{
			name: "hero with YAML block",
			input: `# [hero id="hero-main"] Welcome

This is a description.

` + "```yaml hero\nmore_link: /about\nimage:\n  url: /hero.jpg\n  caption: Hero Image\n```",
			want: []Component{
				&ComponentPageHero{
					ID:          "hero-main",
					Title:       "Welcome",
					Description: "This is a description.",
					MoreLink:    "/about",
					Image: &Image{
						URL:     "/hero.jpg",
						Caption: "Hero Image",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "markdown component",
			input: `# [markdown id="content-1"] 

## This is a heading

Some markdown content here.

- List item 1
- List item 2

More content.`,
			want: []Component{
				&ComponentPageMarkdown{
					ID:       "content-1",
					Markdown: "## This is a heading\n\nSome markdown content here.\n\n- List item 1\n- List item 2\n\nMore content.",
				},
			},
			wantErr: false,
		},
		{
			name: "contact component with YAML",
			input: `# [contact id="contact-section"] Contact Us

` + "```yaml contact\nname: John Doe\nemail: john@example.com\nimage:\n  url: /john.jpg\nsocial_links:\n  - label: Twitter\n    url: https://twitter.com/john\n    external: true\n```",
			want: []Component{
				&ComponentPageContact{
					ID:    "contact-section",
					Name:  "John Doe",
					Email: "john@example.com",
					Image: &Image{
						URL: "/john.jpg",
					},
					SocialLinks: []*SocialLink{
						{
							Label:    "Twitter",
							URL:      "https://twitter.com/john",
							External: true,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "projectSearch component",
			input: `# [projectSearch id="search-1"] 

` + "```yaml projectSearch\nfilter:\n  hiddenCategories:\n    - internal\n  selectedTags:\n    - featured\n```",
			want: []Component{
				&ComponentProjectsProjectSearch{
					ID: "search-1",
					Filter: &ProjectFilter{
						HiddenCategories: []string{"internal"},
						SelectedTags:     []string{"featured"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "highlightReel component",
			input: `# [highlightReel id="highlights"] Featured Projects

` + "```yaml highlightReel\nhighlights:\n  - description: First highlight\n    image:\n      url: /highlight1.jpg\n  - description: Second highlight\n    image:\n      url: /highlight2.jpg\nprojectsToHighlight:\n  - project1\n  - project2\n```",
			want: []Component{
				&ComponentPageHighlightReel{
					ID: "highlights",
					Highlights: []*Highlight{
						{
							Description: "First highlight",
							Image: &Image{
								URL: "/highlight1.jpg",
							},
						},
						{
							Description: "Second highlight",
							Image: &Image{
								URL: "/highlight2.jpg",
							},
						},
					},
					ProjectsToHighlight: []string{"project1", "project2"},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple components",
			input: `# [hero id="hero"] Hero Title

Hero description here.

# [markdown id="content"] Content Section

## Subheading

Some content with **bold** and *italic*.

# [contact id="contact"] Get in Touch

` + "```yaml contact\nname: Jane Doe\nemail: jane@example.com\n```",
			want: []Component{
				&ComponentPageHero{
					ID:          "hero",
					Title:       "Hero Title",
					Description: "Hero description here.",
				},
				&ComponentPageMarkdown{
					ID:       "content",
					Markdown: "## Subheading\n\nSome content with **bold** and *italic*.",
				},
				&ComponentPageContact{
					ID:          "contact",
					Name:        "Jane Doe",
					Email:       "jane@example.com",
					SocialLinks: []*SocialLink{},
				},
			},
			wantErr: false,
		},
		{
			name: "regular headings become markdown components",
			input: `# Regular Heading

Some content under regular heading.

## Another Regular Heading

More content here.

# [hero id="hero"] Hero Component

Hero description.`,
			want: []Component{
				&ComponentPageMarkdown{
					ID:       "markdown-0",
					Markdown: "# Regular Heading\n\nSome content under regular heading.\n\n## Another Regular Heading\n\nMore content here.",
				},
				&ComponentPageHero{
					ID:          "hero",
					Title:       "Hero Component",
					Description: "Hero description.",
				},
			},
			wantErr: false,
		},
		{
			name: "component with multiple attributes",
			input: `# [hero id="main-hero" class="featured"] Welcome`,
			want: []Component{
				&ComponentPageHero{
					ID:    "main-hero",
					Title: "Welcome",
				},
			},
			wantErr: false,
		},
		{
			name:    "empty content",
			input:   "",
			want:    []Component{},
			wantErr: false,
		},
		{
			name: "YAML override of title and description",
			input: `# [hero id="hero"] Original Title

Original description.

` + "```yaml hero\ntitle: YAML Title\ndescription: YAML Description\n```",
			want: []Component{
				&ComponentPageHero{
					ID:          "hero",
					Title:       "YAML Title",
					Description: "YAML Description",
				},
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseMarkdownPage([]byte(tt.input))
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMarkdownPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if len(got) != len(tt.want) {
				t.Errorf("ParseMarkdownPage() returned %d components, want %d", len(got), len(tt.want))
				return
			}
			
			for i, component := range got {
				if !compareComponents(component, tt.want[i]) {
					t.Errorf("ParseMarkdownPage() component %d mismatch:\ngot:  %+v\nwant: %+v", 
						i, component, tt.want[i])
				}
			}
		})
	}
}

// compareComponents compares two components for equality
func compareComponents(a, b Component) bool {
	// Use reflection to compare the concrete types and values
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	
	switch a := a.(type) {
	case *ComponentPageHero:
		b := b.(*ComponentPageHero)
		if a.ID != b.ID || a.Title != b.Title || a.Description != b.Description || a.MoreLink != b.MoreLink {
			return false
		}
		return compareImages(a.Image, b.Image)
		
	case *ComponentPageMarkdown:
		b := b.(*ComponentPageMarkdown)
		return a.ID == b.ID && a.Markdown == b.Markdown
		
	case *ComponentPageContact:
		b := b.(*ComponentPageContact)
		if a.ID != b.ID || a.Name != b.Name || a.Email != b.Email {
			return false
		}
		if !compareImages(a.Image, b.Image) {
			return false
		}
		return compareSocialLinks(a.SocialLinks, b.SocialLinks)
		
	case *ComponentProjectsProjectSearch:
		b := b.(*ComponentProjectsProjectSearch)
		if a.ID != b.ID {
			return false
		}
		return compareProjectFilters(a.Filter, b.Filter)
		
	case *ComponentPageHighlightReel:
		b := b.(*ComponentPageHighlightReel)
		if a.ID != b.ID {
			return false
		}
		if !compareHighlights(a.Highlights, b.Highlights) {
			return false
		}
		if !compareStringSlices(a.ProjectsToHighlight, b.ProjectsToHighlight) {
			return false
		}
		return compareProjectFilters(a.ProjectFilter, b.ProjectFilter)
		
	default:
		return false
	}
}

func compareImages(a, b *Image) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.URL == b.URL && a.Caption == b.Caption && a.Width == b.Width && a.Height == b.Height
}

func compareSocialLinks(a, b []*SocialLink) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Label != b[i].Label || a[i].URL != b[i].URL || 
		   a[i].External != b[i].External || a[i].SVGIcon != b[i].SVGIcon {
			return false
		}
	}
	return true
}

func compareProjectFilters(a, b *ProjectFilter) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return compareStringSlices(a.HiddenCategories, b.HiddenCategories) &&
		compareStringSlices(a.SelectedTags, b.SelectedTags)
}

func compareHighlights(a, b []*Highlight) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Description != b[i].Description {
			return false
		}
		if !compareImages(a[i].Image, b[i].Image) {
			return false
		}
	}
	return true
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParseMarkdownPage_EdgeCases(t *testing.T) {
	parser := NewParser()
	
	tests := []struct {
		name  string
		input string
		check func([]Component) bool
	}{
		{
			name: "component with no ID defaults to empty string",
			input: `# [hero] No ID Hero

Description here.`,
			check: func(components []Component) bool {
				if len(components) != 1 {
					return false
				}
				hero, ok := components[0].(*ComponentPageHero)
				return ok && hero.ID == "" && hero.Title == "No ID Hero"
			},
		},
		{
			name: "trailing whitespace handling",
			input: `# [hero id="hero"]    Title with spaces    

    Description with leading/trailing spaces    

`,
			check: func(components []Component) bool {
				if len(components) != 1 {
					return false
				}
				hero, ok := components[0].(*ComponentPageHero)
				return ok && hero.Title == "Title with spaces" && 
					hero.Description == "Description with leading/trailing spaces"
			},
		},
		{
			name: "YAML block for wrong component type ignored",
			input: `# [hero id="hero"] Hero

` + "```yaml markdown\ncontent: This should be ignored\n```",
			check: func(components []Component) bool {
				if len(components) != 1 {
					return false
				}
				hero, ok := components[0].(*ComponentPageHero)
				return ok && hero.Title == "Hero"
			},
		},
		{
			name: "malformed YAML silently ignored",
			input: `# [hero id="hero"] Title

` + "```yaml hero\nthis is not: valid yaml\n  - mixed indentation\n```",
			check: func(components []Component) bool {
				return len(components) == 1
			},
		},
		{
			name: "nested code blocks preserved in markdown",
			input: "# [markdown id=\"code\"] Code Examples\n\n```javascript\nconst x = 1;\n```\n\nMore text.",
			check: func(components []Component) bool {
				if len(components) != 1 {
					return false
				}
				md, ok := components[0].(*ComponentPageMarkdown)
				return ok && md.ID == "code" && 
					strings.Contains(md.Markdown, "```javascript")
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseMarkdownPage([]byte(tt.input))
			if err != nil {
				t.Errorf("ParseMarkdownPage() unexpected error: %v", err)
				return
			}
			
			if !tt.check(got) {
				t.Errorf("ParseMarkdownPage() check failed for %s", tt.name)
			}
		})
	}
}