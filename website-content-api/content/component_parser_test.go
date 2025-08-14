package content

import (
	"reflect"
	"testing"
)

func TestParseMarkdownImage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Image
	}{
		{
			name:  "valid image",
			input: "![Alt text](/uploads/image.jpg)",
			expected: &Image{
				URL:     "/uploads/image.jpg",
				Caption: "Alt text",
			},
		},
		{
			name:  "image with spaces",
			input: "![A longer caption](/path/to/image.png)",
			expected: &Image{
				URL:     "/path/to/image.png",
				Caption: "A longer caption",
			},
		},
		{
			name:  "empty caption",
			input: "![](/image.jpg)",
			expected: &Image{
				URL:     "/image.jpg",
				Caption: "",
			},
		},
		{
			name:     "not an image",
			input:    "This is not an image",
			expected: nil,
		},
		{
			name:     "malformed image",
			input:    "![Alt text(url)",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMarkdownImage(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseMarkdownImage(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseMarkdownLink(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		expectedURL  string
	}{
		{
			name:         "valid link",
			input:        "[Click here](https://example.com)",
			expectedText: "Click here",
			expectedURL:  "https://example.com",
		},
		{
			name:         "internal link",
			input:        "[About](/about)",
			expectedText: "About",
			expectedURL:  "/about",
		},
		{
			name:         "not a link",
			input:        "Just plain text",
			expectedText: "",
			expectedURL:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, url := parseMarkdownLink(tt.input)
			if text != tt.expectedText || url != tt.expectedURL {
				t.Errorf("parseMarkdownLink(%q) = (%q, %q), want (%q, %q)",
					tt.input, text, url, tt.expectedText, tt.expectedURL)
			}
		})
	}
}

func TestExtractBoldText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold text",
			input:    "**Bold Text**",
			expected: "Bold Text",
		},
		{
			name:     "bold in sentence",
			input:    "This is **important** text",
			expected: "important",
		},
		{
			name:     "no bold text",
			input:    "Regular text",
			expected: "",
		},
		{
			name:     "incomplete bold",
			input:    "**Not closed",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBoldText(tt.input)
			if result != tt.expected {
				t.Errorf("extractBoldText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple list",
			input:    "one,two,three",
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "with spaces",
			input:    "one, two , three",
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single item",
			input:    "single",
			expected: []string{"single"},
		},
		{
			name:     "trailing comma",
			input:    "one,two,",
			expected: []string{"one", "two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitCommaSeparated(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitCommaSeparated(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindSectionContent(t *testing.T) {
	lines := []string{
		"# Main Title",
		"Some intro text",
		"## Section One",
		"Content for section one",
		"More content",
		"## Section Two",
		"Content for section two",
		"### Subsection",
		"Subsection content",
		"## Section Three",
		"Final section",
	}

	tests := []struct {
		name        string
		sectionName string
		expected    []string
	}{
		{
			name:        "find section one",
			sectionName: "Section One",
			expected:    []string{"Content for section one", "More content"},
		},
		{
			name:        "find section two with subsection",
			sectionName: "Section Two",
			expected: []string{
				"Content for section two",
				"### Subsection",
				"Subsection content",
			},
		},
		{
			name:        "find subsection",
			sectionName: "Subsection",
			expected:    []string{"Subsection content"},
		},
		{
			name:        "non-existent section",
			sectionName: "Missing Section",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findSectionContent(lines, tt.sectionName)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("findSectionContent(..., %q) = %v, want %v",
					tt.sectionName, result, tt.expected)
			}
		})
	}
}