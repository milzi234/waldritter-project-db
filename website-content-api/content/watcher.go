package content

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ContentWatcher watches for file changes and reloads content
type ContentWatcher struct {
	vaultPath  string
	loader     *ContentLoader
	store      *ContentStore
	watcher    *fsnotify.Watcher
	debounce   time.Duration
	reloadChan chan struct{}
	stopChan   chan struct{}
	reload     func() error
	mu         sync.Mutex
	running    bool
}

// NewContentWatcher creates a new content watcher
func NewContentWatcher(vaultPath string, loader *ContentLoader, store *ContentStore) (*ContentWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	cw := &ContentWatcher{
		vaultPath:  vaultPath,
		loader:     loader,
		store:      store,
		watcher:    watcher,
		debounce:   100 * time.Millisecond,
		reloadChan: make(chan struct{}, 1),
		stopChan:   make(chan struct{}),
	}

	// Set the default reload function
	cw.reload = cw.defaultReload

	// Add vault directory and subdirectories to watcher
	err = cw.addRecursive(vaultPath)
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to add paths to watcher: %w", err)
	}

	return cw, nil
}

// addRecursive adds a directory and all its subdirectories to the watcher
func (w *ContentWatcher) addRecursive(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories
		if strings.Contains(path, "/.") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			err = w.watcher.Add(path)
			if err != nil {
				log.Printf("Failed to watch directory %s: %v", path, err)
			}
		}

		return nil
	})
}

// Start begins watching for file changes
func (w *ContentWatcher) Start() error {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return fmt.Errorf("watcher already running")
	}
	w.running = true
	w.mu.Unlock()

	go w.watchLoop()
	go w.reloadLoop()

	return nil
}

// Stop stops the watcher
func (w *ContentWatcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return nil
	}

	close(w.stopChan)
	err := w.watcher.Close()
	w.running = false

	return err
}

// watchLoop watches for file system events
func (w *ContentWatcher) watchLoop() {
	var debounceTimer *time.Timer

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Filter out temporary files and non-relevant files
			if !w.shouldReload(event.Name) {
				continue
			}

			// Log the event
			log.Printf("File change detected: %s (%s)", event.Name, event.Op)

			// Debounce rapid changes
			if debounceTimer != nil {
				debounceTimer.Stop()
			}

			debounceTimer = time.AfterFunc(w.debounce, func() {
				select {
				case w.reloadChan <- struct{}{}:
				default: // Channel already has a pending reload
				}
			})

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)

		case <-w.stopChan:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return
		}
	}
}

// reloadLoop processes reload requests
func (w *ContentWatcher) reloadLoop() {
	for {
		select {
		case <-w.reloadChan:
			err := w.reload()
			if err != nil {
				log.Printf("Error reloading content: %v", err)
			}

		case <-w.stopChan:
			return
		}
	}
}

// defaultReload reloads all content from the vault
func (w *ContentWatcher) defaultReload() error {
	log.Println("Reloading content...")

	data, err := w.loader.LoadContent()
	if err != nil {
		return fmt.Errorf("failed to reload content: %w", err)
	}

	w.store.UpdateContent(data.Pages, data.MenuItems, data.Footer)

	stats := w.store.GetStats()
	log.Printf("Content reloaded successfully: %d pages, %d menu items, footer=%v",
		stats["pages_count"], stats["menu_items"], stats["has_footer"])

	return nil
}

// shouldReload determines if a file change should trigger a reload
func (w *ContentWatcher) shouldReload(path string) bool {
	// Ignore hidden files and directories
	if strings.Contains(path, "/.") {
		return false
	}

	// Only reload for markdown files and images
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".jpg", ".jpeg", ".png", ".webp", ".gif", ".svg":
		return true
	default:
		return false
	}
}

// ForceReload manually triggers a content reload
func (w *ContentWatcher) ForceReload() error {
	select {
	case w.reloadChan <- struct{}{}:
		return nil
	default:
		return fmt.Errorf("reload already in progress")
	}
}