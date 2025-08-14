package content

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestVault(t *testing.T) string {
	// Create a temporary directory for the test vault
	tmpDir, err := os.MkdirTemp("", "test-vault-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create directory structure
	dirs := []string{
		"pages",
		"navigation",
		"components",
		"assets/images",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create test page files
	testPages := map[string]string{
		"pages/_home.md": `---
documentId: home
url: /
title: Home
description: Welcome to our site
---

# [hero id="hero-home"] Welcome

This is the homepage

` + "```yaml hero\nimage:\n  url: \"/assets/hero.jpg\"\n  caption: \"Hero image\"\n```" + `

# [markdown id="content-intro"] Introduction

Welcome to our website!

## Additional notes

This is part of the markdown component.`,

		"pages/about.md": `---
documentId: about
url: /about
title: About Us
---

# [markdown id="about-content"] About Us

We are a company.

# [contact id="contact-main"] Contact Us

` + "```yaml contact\nname: \"John Doe\"\nemail: \"john@example.com\"\nsocial_links:\n  - label: \"Twitter\"\n    url: \"https://twitter.com/johndoe\"\n    external: true\n```" + `
`,

		"pages/projects.md": `---
documentId: projects
url: /projects
title: Projects
---

# [projectSearch id="search-projects"] Search Projects

` + "```yaml projectSearch\nfilter:\n  selectedTags:\n    - featured\n```" + `
`,
	}

	for path, content := range testPages {
		fullPath := filepath.Join(tmpDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	// Create test navigation files
	testNav := map[string]string{
		"navigation/main-menu.md": `---
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
      - documentId: menu-projects
        label: "Projects"
        to: "/projects"
        order: 3
---`,

		"navigation/footer.md": `---
type: footer
documentId: main-footer
shortDescriptionTitle: "Test Site"
shortDescriptionText: "This is a test site"
links_1:
  - label: "About"
    url: "/about"
  - label: "Projects"
    url: "/projects"
links_2:
  - label: "Privacy"
    url: "/privacy"
  - label: "Terms"
    url: "/terms"
social_links:
  - label: "Twitter"
    url: "https://twitter.com/test"
    svgIcon: "<svg></svg>"
  - label: "GitHub"
    url: "https://github.com/test"
---`,
	}

	for path, content := range testNav {
		fullPath := filepath.Join(tmpDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	return tmpDir
}

func TestNewContentLoader(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)

	if loader == nil {
		t.Fatal("NewContentLoader() returned nil")
	}

	if loader.vaultPath != tmpDir {
		t.Errorf("NewContentLoader() vaultPath = %v, want %v", loader.vaultPath, tmpDir)
	}

	if loader.parser == nil {
		t.Error("NewContentLoader() parser is nil")
	}
}

func TestContentLoader_LoadContent(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)

	data, err := loader.LoadContent()
	if err != nil {
		t.Fatalf("LoadContent() error = %v", err)
	}

	if data == nil {
		t.Fatal("LoadContent() returned nil data")
	}

	// Check pages
	if len(data.Pages) != 3 {
		t.Errorf("LoadContent() loaded %d pages, want 3", len(data.Pages))
	}

	// Check specific page
	homePage, exists := data.Pages["home"]
	if !exists {
		t.Error("LoadContent() did not load home page")
	} else {
		if homePage.URL != "/" {
			t.Errorf("Home page URL = %v, want /", homePage.URL)
		}
		if len(homePage.Content) != 2 {
			t.Errorf("Home page has %d components, want 2", len(homePage.Content))
		}
	}

	// Check menu items
	if len(data.MenuItems) == 0 {
		t.Error("LoadContent() did not load menu items")
	} else {
		// Check if main-root exists and has items
		foundMain := false
		for _, item := range data.MenuItems {
			if item.Label == "$MAIN" {
				foundMain = true
				if len(item.Items) != 3 {
					t.Errorf("Main menu has %d items, want 3", len(item.Items))
				}
				break
			}
		}
		if !foundMain {
			t.Error("LoadContent() did not load main menu")
		}
	}

	// Check footer
	if data.Footer == nil {
		t.Error("LoadContent() did not load footer")
	} else {
		if data.Footer.ShortDescriptionTitle != "Test Site" {
			t.Errorf("Footer title = %v, want 'Test Site'", data.Footer.ShortDescriptionTitle)
		}
		if len(data.Footer.Links1) != 2 {
			t.Errorf("Footer has %d links_1, want 2", len(data.Footer.Links1))
		}
		if len(data.Footer.SocialLinks) != 2 {
			t.Errorf("Footer has %d social_links, want 2", len(data.Footer.SocialLinks))
		}
	}
}

func TestContentLoader_LoadPage(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)

	pagePath := filepath.Join(tmpDir, "pages", "about.md")
	page, err := loader.loadPage(pagePath)

	if err != nil {
		t.Fatalf("loadPage() error = %v", err)
	}

	if page == nil {
		t.Fatal("loadPage() returned nil")
	}

	if page.DocumentID != "about" {
		t.Errorf("loadPage() DocumentID = %v, want 'about'", page.DocumentID)
	}

	if page.URL != "/about" {
		t.Errorf("loadPage() URL = %v, want '/about'", page.URL)
	}

	if len(page.Content) != 2 {
		t.Errorf("loadPage() loaded %d components, want 2", len(page.Content))
	}

	// Check component types
	if len(page.Content) > 0 {
		switch comp := page.Content[0].(type) {
		case *ComponentPageMarkdown:
			if comp.ID != "about-content" {
				t.Errorf("First component ID = %v, want 'about-content'", comp.ID)
			}
		default:
			t.Errorf("First component type = %T, want *ComponentPageMarkdown", comp)
		}
	}
}

func TestContentLoader_LoadMenu(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)

	menuPath := filepath.Join(tmpDir, "navigation", "main-menu.md")
	menuItems, err := loader.loadMenu(menuPath)

	if err != nil {
		t.Fatalf("loadMenu() error = %v", err)
	}

	if len(menuItems) == 0 {
		t.Fatal("loadMenu() returned empty menu")
	}

	// Check structure
	mainMenu := menuItems[0]
	if mainMenu.Label != "$MAIN" {
		t.Errorf("Main menu label = %v, want '$MAIN'", mainMenu.Label)
	}

	if len(mainMenu.Items) != 3 {
		t.Errorf("Main menu has %d items, want 3", len(mainMenu.Items))
	}

	// Check first item
	if len(mainMenu.Items) > 0 {
		firstItem := mainMenu.Items[0]
		if firstItem.Label != "Home" {
			t.Errorf("First menu item label = %v, want 'Home'", firstItem.Label)
		}
		if firstItem.To != "/" {
			t.Errorf("First menu item to = %v, want '/'", firstItem.To)
		}
	}
}

func TestContentLoader_LoadFooter(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)

	footerPath := filepath.Join(tmpDir, "navigation", "footer.md")
	footer, err := loader.loadFooter(footerPath)

	if err != nil {
		t.Fatalf("loadFooter() error = %v", err)
	}

	if footer == nil {
		t.Fatal("loadFooter() returned nil")
	}

	if footer.ShortDescriptionTitle != "Test Site" {
		t.Errorf("Footer title = %v, want 'Test Site'", footer.ShortDescriptionTitle)
	}

	if len(footer.Links1) != 2 {
		t.Errorf("Footer has %d links_1, want 2", len(footer.Links1))
	}

	if len(footer.Links2) != 2 {
		t.Errorf("Footer has %d links_2, want 2", len(footer.Links2))
	}

	if len(footer.SocialLinks) != 2 {
		t.Errorf("Footer has %d social_links, want 2", len(footer.SocialLinks))
	}

	// Check first social link
	if len(footer.SocialLinks) > 0 {
		firstLink := footer.SocialLinks[0]
		if firstLink.Label != "Twitter" {
			t.Errorf("First social link label = %v, want 'Twitter'", firstLink.Label)
		}
		if firstLink.SVGIcon != "<svg></svg>" {
			t.Errorf("First social link has no SVG icon")
		}
	}
}

func TestContentLoader_NonExistentVault(t *testing.T) {
	loader := NewContentLoader("/non/existent/path")

	data, err := loader.LoadContent()

	// Should not error but return empty data
	if err == nil {
		if len(data.Pages) != 0 {
			t.Error("LoadContent() should return empty pages for non-existent vault")
		}
	}
}

func TestContentLoader_InvalidPageFile(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	// Create an invalid page file
	invalidPath := filepath.Join(tmpDir, "pages", "invalid.md")
	err := os.WriteFile(invalidPath, []byte("invalid: [yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	loader := NewContentLoader(tmpDir)
	data, err := loader.LoadContent()

	// Should continue loading other pages
	if err != nil {
		t.Fatalf("LoadContent() should not fail on invalid page: %v", err)
	}

	// Should have loaded the valid pages
	if len(data.Pages) != 3 {
		t.Errorf("LoadContent() loaded %d pages, want 3 (invalid page should be skipped)", len(data.Pages))
	}
}