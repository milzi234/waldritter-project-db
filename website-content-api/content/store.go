package content

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrPageNotFound is returned when a page is not found
	ErrPageNotFound = errors.New("page not found")
)

// ContentStore manages all content in memory with thread-safe access
type ContentStore struct {
	mu          sync.RWMutex
	pages       map[string]*Page
	menuItems   []*MenuItem
	footer      *Footer
	lastUpdated time.Time
	version     string
}

// NewContentStore creates a new content store
func NewContentStore() *ContentStore {
	return &ContentStore{
		pages:     make(map[string]*Page),
		menuItems: make([]*MenuItem, 0),
		version:   generateVersion(),
	}
}

// generateVersion creates a unique version identifier
func generateVersion() string {
	return uuid.New().String()[:8]
}

// GetPage retrieves a page by its document ID
func (s *ContentStore) GetPage(documentID string) (*Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	page, exists := s.pages[documentID]
	if !exists {
		return nil, ErrPageNotFound
	}

	return page, nil
}

// SetPage stores or updates a page
func (s *ContentStore) SetPage(documentID string, page *Page) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pages[documentID] = page
	s.lastUpdated = time.Now()
}

// GetAllPages returns all pages
func (s *ContentStore) GetAllPages() []*Page {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pages := make([]*Page, 0, len(s.pages))
	for _, page := range s.pages {
		pages = append(pages, page)
	}

	return pages
}

// GetMenuItems returns all menu items
func (s *ContentStore) GetMenuItems() []*MenuItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.menuItems
}

// SetMenuItems updates the menu items
func (s *ContentStore) SetMenuItems(items []*MenuItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.menuItems = items
	s.lastUpdated = time.Now()
}

// GetFooter returns the footer
func (s *ContentStore) GetFooter() *Footer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.footer
}

// SetFooter updates the footer
func (s *ContentStore) SetFooter(footer *Footer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.footer = footer
	s.lastUpdated = time.Now()
}

// UpdateContent atomically updates all content
func (s *ContentStore) UpdateContent(pages map[string]*Page, menu []*MenuItem, footer *Footer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pages = pages
	s.menuItems = menu
	s.footer = footer
	s.lastUpdated = time.Now()
	s.version = generateVersion()
}

// GetLastUpdated returns the last update time
func (s *ContentStore) GetLastUpdated() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastUpdated
}

// GetVersion returns the current content version
func (s *ContentStore) GetVersion() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.version
}

// Clear removes all content from the store
func (s *ContentStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pages = make(map[string]*Page)
	s.menuItems = make([]*MenuItem, 0)
	s.footer = nil
	s.lastUpdated = time.Now()
	s.version = generateVersion()
}

// GetStats returns statistics about the content store
func (s *ContentStore) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"pages_count":   len(s.pages),
		"menu_items":    len(s.menuItems),
		"has_footer":    s.footer != nil,
		"last_updated":  s.lastUpdated,
		"version":       s.version,
	}
}

// ContentData represents all content data for bulk operations
type ContentData struct {
	Pages     map[string]*Page
	MenuItems []*MenuItem
	Footer    *Footer
}