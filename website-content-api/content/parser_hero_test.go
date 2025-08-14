package content

import (
	"reflect"
	"testing"
)

func TestHeroParser(t *testing.T) {
	parser := &HeroParser{}

	tests := []struct {
		name     string
		heading  string
		content  []string
		attrs    map[string]string
		expected *ComponentPageHero
	}{
		{
			name:    "simple hero with all elements",
			heading: "Welcome to Our Site",
			content: []string{
				"This is the description of our amazing website.",
				"",
				"[Learn More](/about)",
				"![Hero Banner](/uploads/hero.jpg)",
			},
			attrs: map[string]string{"id": "main-hero"},
			expected: &ComponentPageHero{
				ID:          "main-hero",
				Title:       "Welcome to Our Site",
				Description: "This is the description of our amazing website.",
				MoreLink:    "/about",
				Image: &Image{
					URL:     "/uploads/hero.jpg",
					Caption: "Hero Banner",
				},
			},
		},
		{
			name:    "hero with multi-line description",
			heading: "Multi-line Hero",
			content: []string{
				"First line of description",
				"Second line of description",
				"Third line of description",
				"",
				"![Image](/img.jpg)",
			},
			attrs: map[string]string{"id": "hero2"},
			expected: &ComponentPageHero{
				ID:          "hero2",
				Title:       "Multi-line Hero",
				Description: "First line of description Second line of description Third line of description",
				Image: &Image{
					URL:     "/img.jpg",
					Caption: "Image",
				},
			},
		},
		{
			name:    "hero with only title",
			heading: "Simple Title",
			content: []string{},
			attrs:   map[string]string{"id": "simple"},
			expected: &ComponentPageHero{
				ID:    "simple",
				Title: "Simple Title",
			},
		},
		{
			name:    "hero without image",
			heading: "No Image Hero",
			content: []string{
				"Just a description here.",
				"[Read More](/more)",
			},
			attrs: map[string]string{"id": "no-img"},
			expected: &ComponentPageHero{
				ID:          "no-img",
				Title:       "No Image Hero",
				Description: "Just a description here.",
				MoreLink:    "/more",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.heading, tt.content, tt.attrs)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			hero, ok := result.(*ComponentPageHero)
			if !ok {
				t.Fatalf("Parse() returned wrong type: %T", result)
			}

			if !reflect.DeepEqual(hero, tt.expected) {
				t.Errorf("Parse() = %+v, want %+v", hero, tt.expected)
			}
		})
	}
}

func TestHeroParserCanParse(t *testing.T) {
	parser := &HeroParser{}

	tests := []struct {
		componentType string
		expected      bool
	}{
		{"hero", true},
		{"Hero", false},
		{"markdown", false},
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