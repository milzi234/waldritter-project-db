package content

import (
	"reflect"
	"testing"
)

func TestParserIntegration(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		content  string
		expected []Component
	}{
		{
			name: "mixed markdown and YAML components",
			content: `# [hero id="main"] Welcome Home
This is our main hero description.
[Learn More](/about)
![Hero Image](/hero.jpg)

# [highlightReel id="highlights"]

## Projects to Highlight
- project-1
- project-2

## Highlights

### First Amazing Feature
![Feature One](/feature1.jpg)

### Second Great Thing
![Feature Two](/feature2.jpg)

# [projectSearch id="search" hide="internal" tags="featured"]

# [contact id="john"] John Doe
john@example.com
![John](/john.jpg)

## Social Links
- [GitHub](https://github.com/john)

# Regular Markdown Section

This is just regular markdown content.
It should be parsed as a markdown component.`,
			expected: []Component{
				&ComponentPageHero{
					ID:          "main",
					Title:       "Welcome Home",
					Description: "This is our main hero description.",
					MoreLink:    "/about",
					Image: &Image{
						URL:     "/hero.jpg",
						Caption: "Hero Image",
					},
				},
				&ComponentPageHighlightReel{
					ID: "highlights",
					ProjectsToHighlight: []string{"project-1", "project-2"},
					Highlights: []*Highlight{
						{
							Description: "First Amazing Feature",
							Image: &Image{
								URL:     "/feature1.jpg",
								Caption: "Feature One",
							},
						},
						{
							Description: "Second Great Thing",
							Image: &Image{
								URL:     "/feature2.jpg",
								Caption: "Feature Two",
							},
						},
					},
				},
				&ComponentProjectsProjectSearch{
					ID: "search",
					Filter: &ProjectFilter{
						HiddenCategories: []string{"internal"},
						SelectedTags:     []string{"featured"},
					},
				},
				&ComponentPageContact{
					ID:    "john",
					Name:  "John Doe",
					Email: "john@example.com",
					Image: &Image{
						URL:     "/john.jpg",
						Caption: "John",
					},
					SocialLinks: []*SocialLink{
						{
							Label:    "GitHub",
							URL:      "https://github.com/john",
							External: true,
							SVGIcon:  "github",
						},
					},
				},
				&ComponentPageMarkdown{
					ID:       "markdown-28", // Line number where the heading appears
					Markdown: "# Regular Markdown Section\n\nThis is just regular markdown content.\nIt should be parsed as a markdown component.",
				},
			},
		},
		{
			name: "YAML override for markdown parsing",
			content: "# [hero id=\"yaml-hero\"] Markdown Title\n\n" +
				"This is markdown description.\n\n" +
				"```yaml hero\n" +
				"title: \"YAML Title Override\"\n" +
				"description: \"YAML Description Override\"\n" +
				"more_link: \"/yaml-link\"\n" +
				"image:\n" +
				"  url: \"/yaml-image.jpg\"\n" +
				"  caption: \"YAML Caption\"\n" +
				"```\n\n" +
				"# [highlightReel id=\"yaml-highlights\"]\n\n" +
				"- **Markdown Highlight** - This would normally be parsed\n" +
				"  ![Markdown Image](/markdown.jpg)\n\n" +
				"```yaml highlightReel\n" +
				"highlights:\n" +
				"  - description: \"YAML Highlight 1\"\n" +
				"    image:\n" +
				"      url: \"/yaml1.jpg\"\n" +
				"      caption: \"YAML 1\"\n" +
				"  - description: \"YAML Highlight 2\"\n" +
				"    image:\n" +
				"      url: \"/yaml2.jpg\"\n" +
				"      caption: \"YAML 2\"\n" +
				"projectsToHighlight:\n" +
				"  - yaml-project-1\n" +
				"  - yaml-project-2\n" +
				"```",
			expected: []Component{
				&ComponentPageHero{
					ID:          "yaml-hero",
					Title:       "YAML Title Override",
					Description: "YAML Description Override",
					MoreLink:    "/yaml-link",
					Image: &Image{
						URL:     "/yaml-image.jpg",
						Caption: "YAML Caption",
					},
				},
				&ComponentPageHighlightReel{
					ID: "yaml-highlights",
					Highlights: []*Highlight{
						{
							Description: "YAML Highlight 1",
							Image: &Image{
								URL:     "/yaml1.jpg",
								Caption: "YAML 1",
							},
						},
						{
							Description: "YAML Highlight 2",
							Image: &Image{
								URL:     "/yaml2.jpg",
								Caption: "YAML 2",
							},
						},
					},
					ProjectsToHighlight: []string{
						"yaml-project-1",
						"yaml-project-2",
					},
				},
			},
		},
		{
			name: "inline attributes parsing",
			content: `# [highlightReel id="inline" projects="p1,p2,p3" hide="private,draft" tags="public"]

- **Quick Feature** - Fast and efficient
  ![Quick](/quick.jpg)

# [projectSearch id="filter" hide="archived,deleted" tags="active,featured"]`,
			expected: []Component{
				&ComponentPageHighlightReel{
					ID:                  "inline",
					ProjectsToHighlight: []string{"p1", "p2", "p3"},
					ProjectFilter: &ProjectFilter{
						HiddenCategories: []string{"private", "draft"},
						SelectedTags:     []string{"public"},
					},
					Highlights: []*Highlight{
						{
							Description: "Fast and efficient",
							Image: &Image{
								URL:     "/quick.jpg",
								Caption: "Quick",
							},
						},
					},
				},
				&ComponentProjectsProjectSearch{
					ID: "filter",
					Filter: &ProjectFilter{
						HiddenCategories: []string{"archived", "deleted"},
						SelectedTags:     []string{"active", "featured"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components, err := parser.ParseMarkdownPage([]byte(tt.content))
			if err != nil {
				t.Fatalf("ParseMarkdownPage() error = %v", err)
			}

			if len(components) != len(tt.expected) {
				t.Logf("Components returned:")
				for i, comp := range components {
					t.Logf("  [%d] %T: %+v", i, comp, comp)
				}
				t.Fatalf("ParseMarkdownPage() returned %d components, want %d", 
					len(components), len(tt.expected))
			}

			for i, component := range components {
				if !reflect.DeepEqual(component, tt.expected[i]) {
					t.Errorf("Component[%d] mismatch:\ngot:  %+v\nwant: %+v",
						i, component, tt.expected[i])
					
					// Type-specific debugging
					switch c := component.(type) {
					case *ComponentPageHero:
						exp := tt.expected[i].(*ComponentPageHero)
						if c.ID != exp.ID {
							t.Errorf("  ID: got %q, want %q", c.ID, exp.ID)
						}
						if c.Title != exp.Title {
							t.Errorf("  Title: got %q, want %q", c.Title, exp.Title)
						}
						if c.Description != exp.Description {
							t.Errorf("  Description: got %q, want %q", c.Description, exp.Description)
						}
						if c.MoreLink != exp.MoreLink {
							t.Errorf("  MoreLink: got %q, want %q", c.MoreLink, exp.MoreLink)
						}
						if !reflect.DeepEqual(c.Image, exp.Image) {
							t.Errorf("  Image: got %+v, want %+v", c.Image, exp.Image)
						}
					case *ComponentPageHighlightReel:
						exp := tt.expected[i].(*ComponentPageHighlightReel)
						if !reflect.DeepEqual(c.ProjectsToHighlight, exp.ProjectsToHighlight) {
							t.Errorf("  ProjectsToHighlight: got %v, want %v", 
								c.ProjectsToHighlight, exp.ProjectsToHighlight)
						}
						if !reflect.DeepEqual(c.ProjectFilter, exp.ProjectFilter) {
							t.Errorf("  ProjectFilter: got %+v, want %+v", 
								c.ProjectFilter, exp.ProjectFilter)
						}
						if len(c.Highlights) != len(exp.Highlights) {
							t.Errorf("  Highlights length: got %d, want %d", 
								len(c.Highlights), len(exp.Highlights))
						}
					case *ComponentProjectsProjectSearch:
						exp := tt.expected[i].(*ComponentProjectsProjectSearch)
						if !reflect.DeepEqual(c.Filter, exp.Filter) {
							t.Errorf("  Filter: got %+v, want %+v", c.Filter, exp.Filter)
						}
					case *ComponentPageContact:
						exp := tt.expected[i].(*ComponentPageContact)
						if c.Name != exp.Name {
							t.Errorf("  Name: got %q, want %q", c.Name, exp.Name)
						}
						if c.Email != exp.Email {
							t.Errorf("  Email: got %q, want %q", c.Email, exp.Email)
						}
					case *ComponentPageMarkdown:
						exp := tt.expected[i].(*ComponentPageMarkdown)
						if c.Markdown != exp.Markdown {
							t.Errorf("  Markdown: got %q, want %q", c.Markdown, exp.Markdown)
						}
					}
				}
			}
		})
	}
}

func TestBackwardCompatibility(t *testing.T) {
	parser := NewParser()

	// Test that existing YAML-only content still works
	yamlOnlyContent := "# [hero id=\"test\"]\n\n" +
		"```yaml hero\n" +
		"title: \"Only YAML Title\"\n" +
		"description: \"Only YAML Description\"\n" +
		"```\n\n" +
		"# [contact id=\"contact\"]\n\n" +
		"```yaml contact\n" +
		"name: \"YAML Contact\"\n" +
		"email: \"yaml@example.com\"\n" +
		"social_links:\n" +
		"  - label: \"Twitter\"\n" +
		"    url: \"https://twitter.com/yaml\"\n" +
		"    external: true\n" +
		"```"

	components, err := parser.ParseMarkdownPage([]byte(yamlOnlyContent))
	if err != nil {
		t.Fatalf("ParseMarkdownPage() error = %v", err)
	}

	if len(components) != 2 {
		t.Fatalf("Expected 2 components, got %d", len(components))
	}

	// Check hero component
	hero, ok := components[0].(*ComponentPageHero)
	if !ok {
		t.Fatalf("First component is not a hero: %T", components[0])
	}
	if hero.Title != "Only YAML Title" {
		t.Errorf("Hero title = %q, want %q", hero.Title, "Only YAML Title")
	}

	// Check contact component
	contact, ok := components[1].(*ComponentPageContact)
	if !ok {
		t.Fatalf("Second component is not a contact: %T", components[1])
	}
	if contact.Name != "YAML Contact" {
		t.Errorf("Contact name = %q, want %q", contact.Name, "YAML Contact")
	}
	if len(contact.SocialLinks) != 1 {
		t.Errorf("Expected 1 social link, got %d", len(contact.SocialLinks))
	}
}