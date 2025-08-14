package content

import (
	"fmt"
	"testing"
	"time"
)

func TestNewContentStore(t *testing.T) {
	store := NewContentStore()
	
	if store == nil {
		t.Fatal("NewContentStore() returned nil")
	}
	
	if store.pages == nil {
		t.Error("NewContentStore() pages map is nil")
	}
	
	if store.menuItems == nil {
		t.Error("NewContentStore() menuItems is nil")
	}
}

func TestContentStore_SetAndGetPage(t *testing.T) {
	store := NewContentStore()
	
	page := &Page{
		DocumentID: "test-page",
		URL:        "/test",
		Content:    []Component{},
		Metadata: &PageMetadata{
			Title:       "Test Page",
			Description: "A test page",
		},
	}
	
	// Set the page
	store.SetPage(page.DocumentID, page)
	
	// Get the page
	retrieved, err := store.GetPage(page.DocumentID)
	if err != nil {
		t.Errorf("GetPage() error = %v", err)
	}
	
	if retrieved == nil {
		t.Fatal("GetPage() returned nil")
	}
	
	if retrieved.DocumentID != page.DocumentID {
		t.Errorf("GetPage() DocumentID = %v, want %v", retrieved.DocumentID, page.DocumentID)
	}
	
	if retrieved.URL != page.URL {
		t.Errorf("GetPage() URL = %v, want %v", retrieved.URL, page.URL)
	}
}

func TestContentStore_GetPage_NotFound(t *testing.T) {
	store := NewContentStore()
	
	_, err := store.GetPage("non-existent")
	if err == nil {
		t.Error("GetPage() expected error for non-existent page")
	}
}

func TestContentStore_GetAllPages(t *testing.T) {
	store := NewContentStore()
	
	// Add multiple pages
	pages := []*Page{
		{DocumentID: "page1", URL: "/page1"},
		{DocumentID: "page2", URL: "/page2"},
		{DocumentID: "page3", URL: "/page3"},
	}
	
	for _, page := range pages {
		store.SetPage(page.DocumentID, page)
	}
	
	// Get all pages
	allPages := store.GetAllPages()
	
	if len(allPages) != len(pages) {
		t.Errorf("GetAllPages() returned %d pages, want %d", len(allPages), len(pages))
	}
	
	// Check that all pages are present
	pageMap := make(map[string]bool)
	for _, page := range allPages {
		pageMap[page.DocumentID] = true
	}
	
	for _, page := range pages {
		if !pageMap[page.DocumentID] {
			t.Errorf("GetAllPages() missing page %s", page.DocumentID)
		}
	}
}

func TestContentStore_SetAndGetMenuItems(t *testing.T) {
	store := NewContentStore()
	
	menuItems := []*MenuItem{
		{
			DocumentID: "menu-home",
			Label:      "Home",
			To:         "/",
			Order:      1,
		},
		{
			DocumentID: "menu-about",
			Label:      "About",
			To:         "/about",
			Order:      2,
		},
	}
	
	store.SetMenuItems(menuItems)
	
	retrieved := store.GetMenuItems()
	
	if len(retrieved) != len(menuItems) {
		t.Errorf("GetMenuItems() returned %d items, want %d", len(retrieved), len(menuItems))
	}
	
	for i, item := range retrieved {
		if item.DocumentID != menuItems[i].DocumentID {
			t.Errorf("GetMenuItems()[%d] DocumentID = %v, want %v", i, item.DocumentID, menuItems[i].DocumentID)
		}
	}
}

func TestContentStore_SetAndGetFooter(t *testing.T) {
	store := NewContentStore()
	
	footer := &Footer{
		ShortDescriptionTitle: "Test Site",
		ShortDescriptionText:  "This is a test site",
		Links1: []*FooterLink{
			{Label: "About", URL: "/about"},
		},
		Links2: []*FooterLink{
			{Label: "Privacy", URL: "/privacy"},
		},
		SocialLinks: []*FooterSocialLink{
			{Label: "Twitter", URL: "https://twitter.com/test"},
		},
	}
	
	store.SetFooter(footer)
	
	retrieved := store.GetFooter()
	
	if retrieved == nil {
		t.Fatal("GetFooter() returned nil")
	}
	
	if retrieved.ShortDescriptionTitle != footer.ShortDescriptionTitle {
		t.Errorf("GetFooter() ShortDescriptionTitle = %v, want %v", retrieved.ShortDescriptionTitle, footer.ShortDescriptionTitle)
	}
	
	if len(retrieved.Links1) != len(footer.Links1) {
		t.Errorf("GetFooter() Links1 length = %v, want %v", len(retrieved.Links1), len(footer.Links1))
	}
}

func TestContentStore_UpdateContent(t *testing.T) {
	store := NewContentStore()
	
	pages := map[string]*Page{
		"page1": {DocumentID: "page1", URL: "/page1"},
		"page2": {DocumentID: "page2", URL: "/page2"},
	}
	
	menuItems := []*MenuItem{
		{DocumentID: "menu1", Label: "Menu 1"},
	}
	
	footer := &Footer{
		ShortDescriptionTitle: "Updated",
		ShortDescriptionText:  "Updated text",
	}
	
	beforeUpdate := time.Now()
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	
	store.UpdateContent(pages, menuItems, footer)
	
	// Check pages were updated
	allPages := store.GetAllPages()
	if len(allPages) != len(pages) {
		t.Errorf("UpdateContent() pages count = %d, want %d", len(allPages), len(pages))
	}
	
	// Check menu items were updated
	retrievedMenu := store.GetMenuItems()
	if len(retrievedMenu) != len(menuItems) {
		t.Errorf("UpdateContent() menu items count = %d, want %d", len(retrievedMenu), len(menuItems))
	}
	
	// Check footer was updated
	retrievedFooter := store.GetFooter()
	if retrievedFooter.ShortDescriptionTitle != footer.ShortDescriptionTitle {
		t.Errorf("UpdateContent() footer title = %v, want %v", retrievedFooter.ShortDescriptionTitle, footer.ShortDescriptionTitle)
	}
	
	// Check last updated time
	if store.GetLastUpdated().Before(beforeUpdate) {
		t.Error("UpdateContent() did not update lastUpdated time")
	}
	
	// Check version was updated
	if store.GetVersion() == "" {
		t.Error("UpdateContent() did not update version")
	}
}

func TestContentStore_ConcurrentAccess(t *testing.T) {
	store := NewContentStore()
	
	// Add initial data
	page := &Page{DocumentID: "test", URL: "/test"}
	store.SetPage("test", page)
	
	done := make(chan bool)
	
	// Concurrent readers
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				store.GetPage("test")
				store.GetAllPages()
				store.GetMenuItems()
				store.GetFooter()
			}
			done <- true
		}()
	}
	
	// Concurrent writers
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 50; j++ {
				p := &Page{
					DocumentID: fmt.Sprintf("page-%d-%d", id, j),
					URL:        fmt.Sprintf("/page-%d-%d", id, j),
				}
				store.SetPage(p.DocumentID, p)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 15; i++ {
		<-done
	}
	
	// Verify data integrity
	allPages := store.GetAllPages()
	if len(allPages) < 1 {
		t.Error("ConcurrentAccess() lost data during concurrent access")
	}
}

func TestContentStore_Clear(t *testing.T) {
	store := NewContentStore()
	
	// Add data
	store.SetPage("page1", &Page{DocumentID: "page1"})
	store.SetMenuItems([]*MenuItem{{DocumentID: "menu1"}})
	store.SetFooter(&Footer{ShortDescriptionTitle: "Test"})
	
	// Clear the store
	store.Clear()
	
	// Verify everything is cleared
	allPages := store.GetAllPages()
	if len(allPages) != 0 {
		t.Errorf("Clear() left %d pages", len(allPages))
	}
	
	menuItems := store.GetMenuItems()
	if len(menuItems) != 0 {
		t.Errorf("Clear() left %d menu items", len(menuItems))
	}
	
	footer := store.GetFooter()
	if footer != nil {
		t.Error("Clear() did not clear footer")
	}
}