package content

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseExampleFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.md")
	
	content := `---
documentId: test-page
url: /test
title: Test Page
---

# [hero id="hero1"] Main Hero Title

Hero description here.

[Learn More](/more)
![Hero Image](/hero.jpg)

# [highlightReel id="highlights"]

## Projects to Highlight
- project-1
- project-2

## Highlights

### Feature One
![Feature 1](/f1.jpg)

### Feature Two
![Feature 2](/f2.jpg)

# [contact id="contact1"] Jane Doe

jane@example.com

## Social Links
- [GitHub](https://github.com/jane)

# Regular Markdown Section

Just regular content here.`

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	// Read and parse the file
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	
	// Extract frontmatter and content
	frontmatter, markdown := ExtractFrontmatter(fileContent)
	
	// Parse metadata
	meta, err := ParsePageMetadata(frontmatter)
	if err != nil {
		t.Fatalf("Failed to parse metadata: %v", err)
	}
	
	// Verify metadata
	if meta.DocumentID != "test-page" {
		t.Errorf("DocumentID = %q, want %q", meta.DocumentID, "test-page")
	}
	if meta.URL != "/test" {
		t.Errorf("URL = %q, want %q", meta.URL, "/test")
	}
	
	// Parse components
	parser := NewParser()
	components, err := parser.ParseMarkdownPage(markdown)
	if err != nil {
		t.Fatalf("Failed to parse components: %v", err)
	}
	
	// Verify we have the expected components
	if len(components) != 4 {
		t.Fatalf("Expected 4 components, got %d", len(components))
	}
	
	// Check component types
	expectedTypes := []string{
		"*content.ComponentPageHero",
		"*content.ComponentPageHighlightReel",
		"*content.ComponentPageContact",
		"*content.ComponentPageMarkdown",
	}
	
	for i, comp := range components {
		typeName := fmt.Sprintf("%T", comp)
		if typeName != expectedTypes[i] {
			t.Errorf("Component[%d] type = %s, want %s", i, typeName, expectedTypes[i])
		}
	}
	
	// Verify hero component
	hero, ok := components[0].(*ComponentPageHero)
	if !ok {
		t.Fatalf("First component is not a hero")
	}
	if hero.Title != "Main Hero Title" {
		t.Errorf("Hero title = %q, want %q", hero.Title, "Main Hero Title")
	}
	if hero.Description != "Hero description here." {
		t.Errorf("Hero description = %q, want %q", hero.Description, "Hero description here.")
	}
	if hero.MoreLink != "/more" {
		t.Errorf("Hero more link = %q, want %q", hero.MoreLink, "/more")
	}
	if hero.Image == nil || hero.Image.URL != "/hero.jpg" {
		t.Errorf("Hero image not parsed correctly")
	}
	
	// Verify highlight reel
	reel, ok := components[1].(*ComponentPageHighlightReel)
	if !ok {
		t.Fatalf("Second component is not a highlight reel")
	}
	if len(reel.ProjectsToHighlight) != 2 {
		t.Errorf("Expected 2 projects to highlight, got %d", len(reel.ProjectsToHighlight))
	}
	if len(reel.Highlights) != 2 {
		t.Errorf("Expected 2 highlights, got %d", len(reel.Highlights))
	}
	
	// Verify contact
	contact, ok := components[2].(*ComponentPageContact)
	if !ok {
		t.Fatalf("Third component is not a contact")
	}
	if contact.Name != "Jane Doe" {
		t.Errorf("Contact name = %q, want %q", contact.Name, "Jane Doe")
	}
	if contact.Email != "jane@example.com" {
		t.Errorf("Contact email = %q, want %q", contact.Email, "jane@example.com")
	}
	if len(contact.SocialLinks) != 1 {
		t.Errorf("Expected 1 social link, got %d", len(contact.SocialLinks))
	}
	
	// Verify markdown section
	md, ok := components[3].(*ComponentPageMarkdown)
	if !ok {
		t.Fatalf("Fourth component is not markdown")
	}
	if md.Markdown != "# Regular Markdown Section\n\nJust regular content here." {
		t.Errorf("Markdown content not as expected")
	}
}

func TestParseComplexHighlightReel(t *testing.T) {
	parser := NewParser()
	
	tests := []struct {
		name    string
		content string
		wantProjects []string
		wantHighlights int
	}{
		{
			name: "inline projects with list highlights",
			content: `# [highlightReel id="test" projects="p1,p2,p3"]

- **Feature A** - Description A
  ![Image A](/a.jpg)
- **Feature B** - Description B
  ![Image B](/b.jpg)`,
			wantProjects: []string{"p1", "p2", "p3"},
			wantHighlights: 2,
		},
		{
			name: "section-based with filter",
			content: `# [highlightReel id="test" hide="draft" tags="public"]

## Projects to Highlight
- project-x
- project-y

## Highlights
### Highlight 1
![Image 1](/1.jpg)`,
			wantProjects: []string{"project-x", "project-y"},
			wantHighlights: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			components, err := parser.ParseMarkdownPage([]byte(tt.content))
			if err != nil {
				t.Fatalf("ParseMarkdownPage() error = %v", err)
			}
			
			if len(components) != 1 {
				t.Fatalf("Expected 1 component, got %d", len(components))
			}
			
			reel, ok := components[0].(*ComponentPageHighlightReel)
			if !ok {
				t.Fatalf("Component is not a highlight reel: %T", components[0])
			}
			
			if len(reel.ProjectsToHighlight) != len(tt.wantProjects) {
				t.Errorf("Projects = %v, want %v", reel.ProjectsToHighlight, tt.wantProjects)
			}
			
			if len(reel.Highlights) != tt.wantHighlights {
				t.Errorf("Got %d highlights, want %d", len(reel.Highlights), tt.wantHighlights)
			}
		})
	}
}