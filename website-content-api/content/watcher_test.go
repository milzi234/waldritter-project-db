package content

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewContentWatcher(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)
	store := NewContentStore()

	watcher, err := NewContentWatcher(tmpDir, loader, store)
	if err != nil {
		t.Fatalf("NewContentWatcher() error = %v", err)
	}

	if watcher == nil {
		t.Fatal("NewContentWatcher() returned nil")
	}

	if watcher.vaultPath != tmpDir {
		t.Errorf("NewContentWatcher() vaultPath = %v, want %v", watcher.vaultPath, tmpDir)
	}

	if watcher.loader != loader {
		t.Error("NewContentWatcher() loader mismatch")
	}

	if watcher.store != store {
		t.Error("NewContentWatcher() store mismatch")
	}

	// Clean up
	err = watcher.Stop()
	if err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}

func TestContentWatcher_FileChange(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)
	store := NewContentStore()

	// Initial load
	data, err := loader.LoadContent()
	if err != nil {
		t.Fatalf("Initial load failed: %v", err)
	}
	store.UpdateContent(data.Pages, data.MenuItems, data.Footer)

	// Create and start watcher
	watcher, err := NewContentWatcher(tmpDir, loader, store)
	if err != nil {
		t.Fatalf("NewContentWatcher() error = %v", err)
	}
	defer watcher.Stop()

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Get initial page count
	initialPages := len(store.GetAllPages())

	// Create a new page file
	newPagePath := filepath.Join(tmpDir, "pages", "new-page.md")
	newPageContent := `---
documentId: new-page
url: /new-page
title: New Page
components:
  - type: markdown
    id: new-content
    content: |
      # New Page
      This is a new page.
---`

	err = os.WriteFile(newPagePath, []byte(newPageContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write new page: %v", err)
	}

	// Wait for the watcher to process the change
	time.Sleep(500 * time.Millisecond)

	// Check that the new page was loaded
	currentPages := len(store.GetAllPages())
	if currentPages <= initialPages {
		t.Errorf("Page count didn't increase after adding file: initial=%d, current=%d", initialPages, currentPages)
	}

	// Check the specific page exists
	newPage, err := store.GetPage("new-page")
	if err != nil {
		t.Errorf("New page not found in store: %v", err)
	}

	if newPage != nil && newPage.URL != "/new-page" {
		t.Errorf("New page URL = %v, want /new-page", newPage.URL)
	}
}

func TestContentWatcher_FileUpdate(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)
	store := NewContentStore()

	// Initial load
	data, err := loader.LoadContent()
	if err != nil {
		t.Fatalf("Initial load failed: %v", err)
	}
	store.UpdateContent(data.Pages, data.MenuItems, data.Footer)

	// Create and start watcher
	watcher, err := NewContentWatcher(tmpDir, loader, store)
	if err != nil {
		t.Fatalf("NewContentWatcher() error = %v", err)
	}
	defer watcher.Stop()

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Get initial about page
	aboutPage, err := store.GetPage("about")
	if err != nil {
		t.Fatalf("About page not found: %v", err)
	}
	initialTitle := aboutPage.Metadata.Title

	// Update the about page
	aboutPath := filepath.Join(tmpDir, "pages", "about.md")
	updatedContent := `---
documentId: about
url: /about
title: Updated About Us
components:
  - type: markdown
    id: about-content
    markdown: |
      # Updated About Us
      
      We are an updated company.
---`

	err = os.WriteFile(aboutPath, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to update about page: %v", err)
	}

	// Wait for the watcher to process the change
	time.Sleep(500 * time.Millisecond)

	// Check that the page was updated
	updatedPage, err := store.GetPage("about")
	if err != nil {
		t.Fatalf("Updated page not found: %v", err)
	}

	if updatedPage.Metadata.Title == initialTitle {
		t.Errorf("Page title was not updated: still %v", initialTitle)
	}

	if updatedPage.Metadata.Title != "Updated About Us" {
		t.Errorf("Page title = %v, want 'Updated About Us'", updatedPage.Metadata.Title)
	}
}

func TestContentWatcher_FileDelete(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)
	store := NewContentStore()

	// Initial load
	data, err := loader.LoadContent()
	if err != nil {
		t.Fatalf("Initial load failed: %v", err)
	}
	store.UpdateContent(data.Pages, data.MenuItems, data.Footer)

	// Create and start watcher
	watcher, err := NewContentWatcher(tmpDir, loader, store)
	if err != nil {
		t.Fatalf("NewContentWatcher() error = %v", err)
	}
	defer watcher.Stop()

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Check projects page exists
	_, err = store.GetPage("projects")
	if err != nil {
		t.Fatalf("Projects page not found initially: %v", err)
	}

	// Delete the projects page
	projectsPath := filepath.Join(tmpDir, "pages", "projects.md")
	err = os.Remove(projectsPath)
	if err != nil {
		t.Fatalf("Failed to delete projects page: %v", err)
	}

	// Wait for the watcher to process the change
	time.Sleep(500 * time.Millisecond)

	// Check that the page was removed
	_, err = store.GetPage("projects")
	if err == nil {
		t.Error("Deleted page still exists in store")
	}
}

func TestContentWatcher_ShouldReload(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)
	store := NewContentStore()

	watcher, err := NewContentWatcher(tmpDir, loader, store)
	if err != nil {
		t.Fatalf("NewContentWatcher() error = %v", err)
	}
	defer watcher.Stop()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"markdown file", "/path/to/file.md", true},
		{"jpg image", "/path/to/image.jpg", true},
		{"png image", "/path/to/image.png", true},
		{"webp image", "/path/to/image.webp", true},
		{"hidden file", "/path/.hidden", false},
		{"hidden directory", "/path/.git/file.md", false},
		{"text file", "/path/to/file.txt", false},
		{"no extension", "/path/to/file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := watcher.shouldReload(tt.path)
			if result != tt.expected {
				t.Errorf("shouldReload(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestContentWatcher_Debounce(t *testing.T) {
	tmpDir := createTestVault(t)
	defer os.RemoveAll(tmpDir)

	loader := NewContentLoader(tmpDir)
	store := NewContentStore()

	// Create watcher with longer debounce for testing
	watcher, err := NewContentWatcher(tmpDir, loader, store)
	if err != nil {
		t.Fatalf("NewContentWatcher() error = %v", err)
	}
	watcher.debounce = 200 * time.Millisecond
	defer watcher.Stop()

	// Track reload count
	reloadCount := 0
	originalReload := watcher.reload
	watcher.reload = func() error {
		reloadCount++
		return originalReload()
	}

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Make multiple rapid changes
	pagePath := filepath.Join(tmpDir, "pages", "test.md")
	for i := 0; i < 5; i++ {
		content := fmt.Sprintf(`---
documentId: test
url: /test
title: Test %d
---`, i)
		err = os.WriteFile(pagePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
		time.Sleep(50 * time.Millisecond) // Less than debounce time
	}

	// Wait for debounce and processing
	time.Sleep(500 * time.Millisecond)

	// Should only reload once due to debouncing
	if reloadCount > 2 {
		t.Errorf("Too many reloads: %d, expected 1 or 2 due to debouncing", reloadCount)
	}
}