package content

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ContentLoader loads content from the Obsidian vault
type ContentLoader struct {
	vaultPath string
	parser    *Parser
}

// NewContentLoader creates a new content loader
func NewContentLoader(vaultPath string) *ContentLoader {
	return &ContentLoader{
		vaultPath: vaultPath,
		parser:    NewParser(),
	}
}

// LoadContent loads all content from the vault
func (l *ContentLoader) LoadContent() (*ContentData, error) {
	data := &ContentData{
		Pages:     make(map[string]*Page),
		MenuItems: make([]*MenuItem, 0),
	}

	// Load pages
	pagesPath := filepath.Join(l.vaultPath, "pages")
	if _, err := os.Stat(pagesPath); err == nil {
		err := filepath.WalkDir(pagesPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v", path, err)
				return nil // Continue walking
			}

			if !d.IsDir() && strings.HasSuffix(path, ".md") {
				page, err := l.loadPage(path)
				if err != nil {
					log.Printf("Error loading page %s: %v", path, err)
					return nil // Continue loading other pages
				}
				if page != nil && page.DocumentID != "" {
					data.Pages[page.DocumentID] = page
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("Error walking pages directory: %v", err)
		}
	}

	// Load navigation
	menuPath := filepath.Join(l.vaultPath, "navigation", "main-menu.md")
	if _, err := os.Stat(menuPath); err == nil {
		menuItems, err := l.loadMenu(menuPath)
		if err != nil {
			log.Printf("Error loading menu: %v", err)
		} else {
			data.MenuItems = menuItems
		}
	}

	// Load footer
	footerPath := filepath.Join(l.vaultPath, "navigation", "footer.md")
	if _, err := os.Stat(footerPath); err == nil {
		footer, err := l.loadFooter(footerPath)
		if err != nil {
			log.Printf("Error loading footer: %v", err)
		} else {
			data.Footer = footer
		}
	}

	return data, nil
}

// loadPage loads a single page from a file
func (l *ContentLoader) loadPage(path string) (*Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read page file: %w", err)
	}

	// Parse frontmatter and content
	frontmatter, markdownBody := ExtractFrontmatter(content)

	// Parse metadata
	pageMeta, err := ParsePageMetadata(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page metadata: %w", err)
	}

	// Skip pages without documentId
	if pageMeta.DocumentID == "" {
		return nil, fmt.Errorf("page missing documentId: %s", path)
	}

	// Parse components from markdown body using the new parser
	components, err := l.parser.ParseMarkdownPage(markdownBody)
	if err != nil {
		log.Printf("Error parsing markdown components in %s: %v", path, err)
		// Return page with empty components rather than failing completely
		components = []Component{}
	}

	// Get file modification time
	fileInfo, err := os.Stat(path)
	var lastModified *time.Time
	if err == nil {
		modTime := fileInfo.ModTime()
		lastModified = &modTime
	}

	// Use lastModified from metadata if available, otherwise use file mod time
	if pageMeta.LastModified != nil {
		lastModified = pageMeta.LastModified
	}

	return &Page{
		DocumentID: pageMeta.DocumentID,
		URL:        pageMeta.URL,
		Content:    components,
		Metadata: &PageMetadata{
			Title:        pageMeta.Title,
			Description:  pageMeta.Description,
			LastModified: lastModified,
			Author:       pageMeta.Author,
		},
	}, nil
}

// loadMenu loads the menu structure from a file
func (l *ContentLoader) loadMenu(path string) ([]*MenuItem, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read menu file: %w", err)
	}

	// Extract frontmatter and markdown content
	frontmatter, markdownContent := ExtractFrontmatter(content)

	// Check if YAML frontmatter contains menuItems - if so, use YAML parser
	if len(frontmatter) > 0 {
		// Try to parse metadata to check if it has menuItems
		menuMeta, err := ParseMenuMetadata(frontmatter)
		if err == nil && len(menuMeta.MenuItems) > 0 {
			log.Printf("Menu loaded using YAML parser from %s", path)
			// Convert to GraphQL types
			return ConvertMenuItems(menuMeta.MenuItems), nil
		}
	}

	// If no YAML menuItems found, try markdown parsing
	if len(markdownContent) > 0 {
		// Attempt to parse as markdown
		menuItems, err := ParseMarkdownMenu(markdownContent)
		if err == nil && len(menuItems) > 0 {
			// Check if the parsed menu has actual items (not just empty $MAIN root)
			if len(menuItems) > 0 && len(menuItems[0].Items) > 0 {
				log.Printf("Menu loaded using markdown parser from %s", path)
				return menuItems, nil
			}
		}
		// If markdown parsing fails, log the error
		if err != nil {
			log.Printf("Markdown menu parsing failed: %v", err)
		}
	}

	// If both failed, return an error
	return nil, fmt.Errorf("no valid menu content found in %s", path)
}

// loadFooter loads the footer content from a file
func (l *ContentLoader) loadFooter(path string) (*Footer, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read footer file: %w", err)
	}

	// Parse frontmatter
	frontmatter, _ := ExtractFrontmatter(content)

	// Parse metadata
	footerMeta, err := ParseFooterMetadata(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("failed to parse footer metadata: %w", err)
	}

	// Convert to GraphQL types
	return ConvertFooter(footerMeta), nil
}

// ValidateContent checks if all required content is present
func (l *ContentLoader) ValidateContent(data *ContentData) []string {
	var errors []string

	// Check for home page
	if _, hasHome := data.Pages["home"]; !hasHome {
		if _, hasIndexHome := data.Pages["_home"]; !hasIndexHome {
			errors = append(errors, "Missing home page (documentId: 'home' or '_home')")
		}
	}

	// Check for empty pages
	if len(data.Pages) == 0 {
		errors = append(errors, "No pages found")
	}

	// Check for menu
	if len(data.MenuItems) == 0 {
		errors = append(errors, "No menu items found")
	}

	// Check for footer
	if data.Footer == nil {
		errors = append(errors, "No footer found")
	}

	// Check page URLs
	for id, page := range data.Pages {
		if page.URL == "" {
			errors = append(errors, fmt.Sprintf("Page %s missing URL", id))
		}
		if page.DocumentID == "" {
			errors = append(errors, fmt.Sprintf("Page %s missing documentId", id))
		}
	}

	return errors
}

// GetAssetPath returns the full path to an asset file
func (l *ContentLoader) GetAssetPath(relativePath string) string {
	// Remove leading slash if present
	relativePath = strings.TrimPrefix(relativePath, "/")
	
	// If it starts with "assets/", use it directly
	if strings.HasPrefix(relativePath, "assets/") {
		return filepath.Join(l.vaultPath, relativePath)
	}
	
	// Otherwise, assume it's under assets
	return filepath.Join(l.vaultPath, "assets", relativePath)
}

// AssetExists checks if an asset file exists
func (l *ContentLoader) AssetExists(relativePath string) bool {
	assetPath := l.GetAssetPath(relativePath)
	_, err := os.Stat(assetPath)
	return err == nil
}