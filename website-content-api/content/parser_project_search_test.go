package content

import (
	"reflect"
	"testing"
)

func TestProjectSearchParser(t *testing.T) {
	parser := &ProjectSearchParser{}

	tests := []struct {
		name     string
		heading  string
		content  []string
		attrs    map[string]string
		expected *ComponentProjectsProjectSearch
	}{
		{
			name:    "inline attributes",
			heading: "",
			content: []string{},
			attrs: map[string]string{
				"id":   "filter1",
				"hide": "internal,archive",
				"tags": "featured,current",
			},
			expected: &ComponentProjectsProjectSearch{
				ID: "filter1",
				Filter: &ProjectFilter{
					HiddenCategories: []string{"internal", "archive"},
					SelectedTags:     []string{"featured", "current"},
				},
			},
		},
		{
			name:    "section-based format",
			heading: "",
			content: []string{
				"## Filter Settings",
				"- Hide categories: internal, archive",
				"- Show tags: featured, current",
			},
			attrs: map[string]string{"id": "filter2"},
			expected: &ComponentProjectsProjectSearch{
				ID: "filter2",
				Filter: &ProjectFilter{
					HiddenCategories: []string{"internal", "archive"},
					SelectedTags:     []string{"featured", "current"},
				},
			},
		},
		{
			name:    "list without section",
			heading: "",
			content: []string{
				"- Hide categories: test1, test2",
				"- Selected tags: tag1, tag2, tag3",
			},
			attrs: map[string]string{"id": "filter3"},
			expected: &ComponentProjectsProjectSearch{
				ID: "filter3",
				Filter: &ProjectFilter{
					HiddenCategories: []string{"test1", "test2"},
					SelectedTags:     []string{"tag1", "tag2", "tag3"},
				},
			},
		},
		{
			name:    "non-list format",
			heading: "",
			content: []string{
				"Hide categories: private",
				"Show tags: public, visible",
			},
			attrs: map[string]string{"id": "filter4"},
			expected: &ComponentProjectsProjectSearch{
				ID: "filter4",
				Filter: &ProjectFilter{
					HiddenCategories: []string{"private"},
					SelectedTags:     []string{"public", "visible"},
				},
			},
		},
		{
			name:    "empty filter",
			heading: "",
			content: []string{},
			attrs:   map[string]string{"id": "empty-filter"},
			expected: &ComponentProjectsProjectSearch{
				ID: "empty-filter",
			},
		},
		{
			name:    "mixed inline and content",
			heading: "",
			content: []string{
				"- Hide categories: extra1, extra2",
			},
			attrs: map[string]string{
				"id":   "mixed",
				"tags": "inline1,inline2",
			},
			expected: &ComponentProjectsProjectSearch{
				ID: "mixed",
				Filter: &ProjectFilter{
					HiddenCategories: []string{"extra1", "extra2"},
					SelectedTags:     []string{"inline1", "inline2"},
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

			search, ok := result.(*ComponentProjectsProjectSearch)
			if !ok {
				t.Fatalf("Parse() returned wrong type: %T", result)
			}

			if !reflect.DeepEqual(search, tt.expected) {
				t.Errorf("Parse() = %+v, want %+v", search, tt.expected)
				
				// Detailed comparison for debugging
				if search.Filter == nil && tt.expected.Filter == nil {
					// Both nil, OK
				} else if search.Filter == nil || tt.expected.Filter == nil {
					t.Errorf("Filter nil mismatch: got %v, want %v",
						search.Filter, tt.expected.Filter)
				} else {
					if !reflect.DeepEqual(search.Filter.HiddenCategories, tt.expected.Filter.HiddenCategories) {
						t.Errorf("HiddenCategories mismatch: got %v, want %v",
							search.Filter.HiddenCategories, tt.expected.Filter.HiddenCategories)
					}
					if !reflect.DeepEqual(search.Filter.SelectedTags, tt.expected.Filter.SelectedTags) {
						t.Errorf("SelectedTags mismatch: got %v, want %v",
							search.Filter.SelectedTags, tt.expected.Filter.SelectedTags)
					}
				}
			}
		})
	}
}

func TestProjectSearchParserCanParse(t *testing.T) {
	parser := &ProjectSearchParser{}

	tests := []struct {
		componentType string
		expected      bool
	}{
		{"projectSearch", true},
		{"ProjectSearch", false},
		{"search", false},
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