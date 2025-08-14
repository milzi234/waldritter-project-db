package content

import (
	"reflect"
	"testing"
)

func TestHighlightReelParser(t *testing.T) {
	parser := &HighlightReelParser{}

	tests := []struct {
		name     string
		heading  string
		content  []string
		attrs    map[string]string
		expected *ComponentPageHighlightReel
	}{
		{
			name:    "section-based format",
			heading: "",
			content: []string{
				"## Projects to Highlight",
				"- larp-fuer-demokratie",
				"- sommerfreizeit-2024",
				"",
				"## Highlights",
				"",
				"### LARP für Demokratie - Spielerisch Demokratie verstehen",
				"![Demokratiebildung durch Rollenspiel](/uploads/campfire.webp)",
				"",
				"### 10 Ortsgruppen bundesweit",
				"![Gemeinsam Geschichten erleben](/uploads/anna.webp)",
			},
			attrs: map[string]string{"id": "featured"},
			expected: &ComponentPageHighlightReel{
				ID: "featured",
				ProjectsToHighlight: []string{
					"larp-fuer-demokratie",
					"sommerfreizeit-2024",
				},
				Highlights: []*Highlight{
					{
						Description: "LARP für Demokratie - Spielerisch Demokratie verstehen",
						Image: &Image{
							URL:     "/uploads/campfire.webp",
							Caption: "Demokratiebildung durch Rollenspiel",
						},
					},
					{
						Description: "10 Ortsgruppen bundesweit",
						Image: &Image{
							URL:     "/uploads/anna.webp",
							Caption: "Gemeinsam Geschichten erleben",
						},
					},
				},
			},
		},
		{
			name:    "inline attributes with list format",
			heading: "",
			content: []string{
				"- **LARP für Demokratie** - Spielerisch Demokratie verstehen",
				"  ![Demokratiebildung durch Rollenspiel](/uploads/campfire.webp)",
				"",
				"- **10 Ortsgruppen** - Finde deine Abenteuergruppe",
				"  ![Gemeinsam erleben](/uploads/anna.webp)",
			},
			attrs: map[string]string{
				"id":       "featured2",
				"projects": "project1,project2",
			},
			expected: &ComponentPageHighlightReel{
				ID: "featured2",
				ProjectsToHighlight: []string{"project1", "project2"},
				Highlights: []*Highlight{
					{
						Description: "Spielerisch Demokratie verstehen",
						Image: &Image{
							URL:     "/uploads/campfire.webp",
							Caption: "Demokratiebildung durch Rollenspiel",
						},
					},
					{
						Description: "Finde deine Abenteuergruppe",
						Image: &Image{
							URL:     "/uploads/anna.webp",
							Caption: "Gemeinsam erleben",
						},
					},
				},
			},
		},
		{
			name:    "with filter settings",
			heading: "",
			content: []string{
				"## Highlights",
				"### First Highlight",
				"![Image 1](/img1.jpg)",
			},
			attrs: map[string]string{
				"id":   "filtered",
				"hide": "internal,archive",
				"tags": "featured,current",
			},
			expected: &ComponentPageHighlightReel{
				ID: "filtered",
				Highlights: []*Highlight{
					{
						Description: "First Highlight",
						Image: &Image{
							URL:     "/img1.jpg",
							Caption: "Image 1",
						},
					},
				},
				ProjectFilter: &ProjectFilter{
					HiddenCategories: []string{"internal", "archive"},
					SelectedTags:     []string{"featured", "current"},
				},
			},
		},
		{
			name:    "simple list without bold text",
			heading: "",
			content: []string{
				"- First highlight description",
				"  ![First image](/first.jpg)",
				"- Second highlight description",
				"  ![Second image](/second.jpg)",
			},
			attrs: map[string]string{"id": "simple-list"},
			expected: &ComponentPageHighlightReel{
				ID: "simple-list",
				Highlights: []*Highlight{
					{
						Description: "First highlight description",
						Image: &Image{
							URL:     "/first.jpg",
							Caption: "First image",
						},
					},
					{
						Description: "Second highlight description",
						Image: &Image{
							URL:     "/second.jpg",
							Caption: "Second image",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.heading, tt.content, tt.attrs)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			reel, ok := result.(*ComponentPageHighlightReel)
			if !ok {
				t.Fatalf("Parse() returned wrong type: %T", result)
			}

			if !reflect.DeepEqual(reel, tt.expected) {
				t.Errorf("Parse() = %+v, want %+v", reel, tt.expected)
				
				// Detailed comparison for debugging
				if !reflect.DeepEqual(reel.ProjectsToHighlight, tt.expected.ProjectsToHighlight) {
					t.Errorf("ProjectsToHighlight mismatch: got %v, want %v",
						reel.ProjectsToHighlight, tt.expected.ProjectsToHighlight)
				}
				if !reflect.DeepEqual(reel.ProjectFilter, tt.expected.ProjectFilter) {
					t.Errorf("ProjectFilter mismatch: got %+v, want %+v",
						reel.ProjectFilter, tt.expected.ProjectFilter)
				}
				if len(reel.Highlights) != len(tt.expected.Highlights) {
					t.Errorf("Highlights length mismatch: got %d, want %d",
						len(reel.Highlights), len(tt.expected.Highlights))
				} else {
					for i, h := range reel.Highlights {
						if !reflect.DeepEqual(h, tt.expected.Highlights[i]) {
							t.Errorf("Highlight[%d] mismatch: got %+v, want %+v",
								i, h, tt.expected.Highlights[i])
						}
					}
				}
			}
		})
	}
}

func TestHighlightReelParserCanParse(t *testing.T) {
	parser := &HighlightReelParser{}

	tests := []struct {
		componentType string
		expected      bool
	}{
		{"highlightReel", true},
		{"highlight", false},
		{"reel", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.componentType, func(t *testing.T) {
			result := parser.CanParse(tt.componentType)
			if result != tt.expected {
				t.Errorf("CanParse(%q) = %v, want %v", tt.componentType, result, tt.expected)
			}
		})
	}
}