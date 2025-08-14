package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPage_WithMarkdownComponents(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	pagesDir := filepath.Join(tempDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	
	// Test cases with different page formats
	tests := []struct {
		name        string
		filename    string
		content     string
		wantPageID  string
		wantURL     string
		checkComps  func([]Component) bool
	}{
		{
			name:     "page with markdown components",
			filename: "home.md",
			content: `---
documentId: home
url: /
title: Home Page
description: Welcome to our site
author: Test Author
---

# [hero id="hero-main"] Welcome to Our Site

This is the main hero section of our homepage.

` + "```yaml hero\nmore_link: /about\nimage:\n  url: /hero-image.jpg\n```" + `

# [markdown id="intro"] Introduction

## About Us

We are a company dedicated to excellence.

* First point
* Second point
* Third point

# [contact id="contact-form"] Get in Touch

` + "```yaml contact\nname: Support Team\nemail: support@example.com\n```",
			wantPageID: "home",
			wantURL:    "/",
			checkComps: func(comps []Component) bool {
				if len(comps) != 3 {
					return false
				}
				
				// Check hero component
				hero, ok := comps[0].(*ComponentPageHero)
				if !ok || hero.ID != "hero-main" || 
				   hero.Title != "Welcome to Our Site" ||
				   hero.MoreLink != "/about" ||
				   hero.Image == nil || hero.Image.URL != "/hero-image.jpg" {
					return false
				}
				
				// Check markdown component
				md, ok := comps[1].(*ComponentPageMarkdown)
				if !ok || md.ID != "intro" || 
				   !strings.Contains(md.Markdown, "## About Us") {
					return false
				}
				
				// Check contact component
				contact, ok := comps[2].(*ComponentPageContact)
				if !ok || contact.ID != "contact-form" ||
				   contact.Name != "Support Team" ||
				   contact.Email != "support@example.com" {
					return false
				}
				
				return true
			},
		},
		{
			name:     "page with project search and highlight reel",
			filename: "projects.md",
			content: `---
documentId: projects
url: /projects
title: Projects
---

# [projectSearch id="search"] Find Projects

` + "```yaml projectSearch\nfilter:\n  selectedTags:\n    - featured\n    - public\n```" + `

# [highlightReel id="highlights"] Featured Work

` + "```yaml highlightReel\nprojectsToHighlight:\n  - project-1\n  - project-2\nhighlights:\n  - description: First highlight\n    image:\n      url: /highlight1.jpg\n```",
			wantPageID: "projects",
			wantURL:    "/projects",
			checkComps: func(comps []Component) bool {
				if len(comps) != 2 {
					return false
				}
				
				// Check project search
				search, ok := comps[0].(*ComponentProjectsProjectSearch)
				if !ok || search.ID != "search" ||
				   search.Filter == nil ||
				   len(search.Filter.SelectedTags) != 2 {
					return false
				}
				
				// Check highlight reel
				highlights, ok := comps[1].(*ComponentPageHighlightReel)
				if !ok || highlights.ID != "highlights" ||
				   len(highlights.ProjectsToHighlight) != 2 ||
				   len(highlights.Highlights) != 1 {
					return false
				}
				
				return true
			},
		},
		{
			name:     "page with regular markdown that becomes components",
			filename: "about.md",
			content: `---
documentId: about
url: /about
title: About
---

# About Our Company

We started in 2020 with a simple mission.

## Our Mission

To create amazing experiences.

## Our Team

We have a diverse team of professionals.

# [contact id="team-contact"] Contact Our Team

` + "```yaml contact\nemail: team@example.com\n```",
			wantPageID: "about",
			wantURL:    "/about",
			checkComps: func(comps []Component) bool {
				if len(comps) != 2 {
					return false
				}
				
				// First component should be markdown with all regular content
				md, ok := comps[0].(*ComponentPageMarkdown)
				if !ok || !strings.Contains(md.Markdown, "# About Our Company") ||
				   !strings.Contains(md.Markdown, "## Our Mission") ||
				   !strings.Contains(md.Markdown, "## Our Team") {
					return false
				}
				
				// Second component should be contact
				contact, ok := comps[1].(*ComponentPageContact)
				if !ok || contact.ID != "team-contact" ||
				   contact.Email != "team@example.com" {
					return false
				}
				
				return true
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write test file
			pagePath := filepath.Join(pagesDir, tt.filename)
			err := os.WriteFile(pagePath, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}
			
			// Create loader and load the page
			loader := NewContentLoader(tempDir)
			page, err := loader.loadPage(pagePath)
			
			if err != nil {
				t.Errorf("loadPage() error = %v", err)
				return
			}
			
			if page.DocumentID != tt.wantPageID {
				t.Errorf("loadPage() DocumentID = %v, want %v", page.DocumentID, tt.wantPageID)
			}
			
			if page.URL != tt.wantURL {
				t.Errorf("loadPage() URL = %v, want %v", page.URL, tt.wantURL)
			}
			
			if !tt.checkComps(page.Content) {
				t.Errorf("loadPage() components check failed for %s", tt.name)
			}
		})
	}
}

func TestLoadPage_BackwardCompatibility(t *testing.T) {
	// This test ensures that old format pages still work
	tempDir := t.TempDir()
	pagesDir := filepath.Join(tempDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	
	// Old format with components in frontmatter (should now parse from body)
	oldFormatContent := `---
documentId: legacy-page
url: /legacy
title: Legacy Page
---

# This is a regular heading

Some content that should become a markdown component.

## Another heading

More content here.`
	
	pagePath := filepath.Join(pagesDir, "legacy.md")
	err := os.WriteFile(pagePath, []byte(oldFormatContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	loader := NewContentLoader(tempDir)
	page, err := loader.loadPage(pagePath)
	
	if err != nil {
		t.Errorf("loadPage() error = %v", err)
		return
	}
	
	if page.DocumentID != "legacy-page" {
		t.Errorf("loadPage() DocumentID = %v, want legacy-page", page.DocumentID)
	}
	
	// Should have created a markdown component from the content
	if len(page.Content) != 1 {
		t.Errorf("loadPage() expected 1 component, got %d", len(page.Content))
		return
	}
	
	md, ok := page.Content[0].(*ComponentPageMarkdown)
	if !ok {
		t.Errorf("loadPage() expected ComponentPageMarkdown, got %T", page.Content[0])
		return
	}
	
	if !strings.Contains(md.Markdown, "# This is a regular heading") {
		t.Errorf("loadPage() markdown content missing expected heading")
	}
}

func TestLoadContent_Integration(t *testing.T) {
	// Integration test for the full content loading process
	tempDir := t.TempDir()
	
	// Create directory structure
	pagesDir := filepath.Join(tempDir, "pages")
	navDir := filepath.Join(tempDir, "navigation")
	os.MkdirAll(pagesDir, 0755)
	os.MkdirAll(navDir, 0755)
	
	// Create a test page with new format
	homeContent := `---
documentId: home
url: /
title: Home
---

# [hero id="main-hero"] Welcome Home

The best place to start.

# [markdown id="content"] Content Section

## Important Information

This is important content that needs to be displayed.`
	
	err := os.WriteFile(filepath.Join(pagesDir, "home.md"), []byte(homeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write home page: %v", err)
	}
	
	// Create another page
	aboutContent := `---
documentId: about
url: /about
title: About
---

# [markdown id="about-content"] About Us

We are a great company.`
	
	err = os.WriteFile(filepath.Join(pagesDir, "about.md"), []byte(aboutContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write about page: %v", err)
	}
	
	// Create menu file
	menuContent := `---
type: menu
documentId: main-menu
menuItems:
  - documentId: home
    label: Home
    to: /
  - documentId: about
    label: About
    to: /about
---`
	
	err = os.WriteFile(filepath.Join(navDir, "main-menu.md"), []byte(menuContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write menu file: %v", err)
	}
	
	// Create footer file
	footerContent := `---
type: footer
documentId: main-footer
shortDescriptionTitle: Our Company
shortDescriptionText: Building great things
links_1:
  - label: Privacy
    url: /privacy
links_2:
  - label: Contact
    url: /contact
---`
	
	err = os.WriteFile(filepath.Join(navDir, "footer.md"), []byte(footerContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write footer file: %v", err)
	}
	
	// Load all content
	loader := NewContentLoader(tempDir)
	data, err := loader.LoadContent()
	
	if err != nil {
		t.Errorf("LoadContent() error = %v", err)
		return
	}
	
	// Check pages loaded
	if len(data.Pages) != 2 {
		t.Errorf("LoadContent() loaded %d pages, want 2", len(data.Pages))
	}
	
	// Check home page
	homePage, ok := data.Pages["home"]
	if !ok {
		t.Error("LoadContent() home page not found")
	} else {
		if len(homePage.Content) != 2 {
			t.Errorf("LoadContent() home page has %d components, want 2", len(homePage.Content))
		}
		
		// Check first component is hero
		if hero, ok := homePage.Content[0].(*ComponentPageHero); ok {
			if hero.ID != "main-hero" || hero.Title != "Welcome Home" {
				t.Errorf("LoadContent() home hero component incorrect")
			}
		} else {
			t.Errorf("LoadContent() home first component is not hero")
		}
	}
	
	// Check about page
	aboutPage, ok := data.Pages["about"]
	if !ok {
		t.Error("LoadContent() about page not found")
	} else {
		if len(aboutPage.Content) != 1 {
			t.Errorf("LoadContent() about page has %d components, want 1", len(aboutPage.Content))
		}
	}
	
	// Check menu loaded
	if len(data.MenuItems) == 0 {
		t.Error("LoadContent() no menu items loaded")
	}
	
	// Check footer loaded
	if data.Footer == nil {
		t.Error("LoadContent() footer not loaded")
	} else {
		if data.Footer.ShortDescriptionTitle != "Our Company" {
			t.Errorf("LoadContent() footer title = %v, want 'Our Company'", 
				data.Footer.ShortDescriptionTitle)
		}
	}
}