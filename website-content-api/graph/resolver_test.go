package graph_test

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/waldritter/website-content-api/content"
	"github.com/waldritter/website-content-api/graph"
	"github.com/waldritter/website-content-api/graph/generated"
)

func setupTestServer(t *testing.T) *client.Client {
	// Create test content
	store := content.NewContentStore()
	
	// Add test pages
	testPages := map[string]*content.Page{
		"home": {
			DocumentID: "home",
			URL:        "/",
			Content: []content.Component{
				&content.ComponentPageHero{
					ID:          "hero-1",
					Title:       "Welcome",
					Description: "Test description",
				},
				&content.ComponentPageMarkdown{
					ID:       "markdown-1",
					Markdown: "# Test Content\n\nThis is test content.",
				},
			},
			Metadata: &content.PageMetadata{
				Title:       "Home Page",
				Description: "Home page description",
			},
		},
		"about": {
			DocumentID: "about",
			URL:        "/about",
			Content:    []content.Component{},
			Metadata: &content.PageMetadata{
				Title: "About Page",
			},
		},
	}
	
	// Add test menu
	testMenu := []*content.MenuItem{
		{
			DocumentID: "main-root",
			Label:      "$MAIN",
			Items: []*content.MenuItem{
				{DocumentID: "menu-home", Label: "Home", To: "/", Order: 1},
				{DocumentID: "menu-about", Label: "About", To: "/about", Order: 2},
			},
		},
	}
	
	// Add test footer
	testFooter := &content.Footer{
		ShortDescriptionTitle: "Test Site",
		ShortDescriptionText:  "Test description",
		Links1: []*content.FooterLink{
			{Label: "About", URL: "/about"},
		},
		Links2: []*content.FooterLink{
			{Label: "Privacy", URL: "/privacy"},
		},
		SocialLinks: []*content.FooterSocialLink{
			{Label: "Twitter", URL: "https://twitter.com/test"},
		},
	}
	
	store.UpdateContent(testPages, testMenu, testFooter)
	
	// Create resolver and server
	resolver := &graph.Resolver{
		Store: store,
	}
	
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	
	return client.New(srv)
}

func TestPagesQuery(t *testing.T) {
	c := setupTestServer(t)
	
	var resp struct {
		Pages []struct {
			DocumentID string `json:"documentId"`
			URL        string `json:"url"`
		} `json:"pages"`
	}
	
	query := `
		query {
			pages {
				documentId
				url
			}
		}
	`
	
	err := c.Post(query, &resp)
	if err != nil {
		t.Fatalf("Pages query failed: %v", err)
	}
	
	if len(resp.Pages) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(resp.Pages))
	}
	
	// Check for home page
	foundHome := false
	for _, page := range resp.Pages {
		if page.DocumentID == "home" && page.URL == "/" {
			foundHome = true
			break
		}
	}
	
	if !foundHome {
		t.Error("Home page not found in pages query")
	}
}

func TestPageQuery(t *testing.T) {
	c := setupTestServer(t)
	
	var resp struct {
		Page struct {
			DocumentID string `json:"documentId"`
			URL        string `json:"url"`
			Content    []map[string]interface{} `json:"content"`
			Metadata   struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"metadata"`
		} `json:"page"`
	}
	
	query := `
		query LoadPage($documentId: ID!) {
			page(documentId: $documentId) {
				documentId
				url
				content {
					__typename
					... on ComponentPageHero {
						id
						title
						description
					}
					... on ComponentPageMarkdown {
						id
						markdown
					}
				}
				metadata {
					title
					description
				}
			}
		}
	`
	
	err := c.Post(query, &resp, client.Var("documentId", "home"))
	if err != nil {
		t.Fatalf("Page query failed: %v", err)
	}
	
	if resp.Page.DocumentID != "home" {
		t.Errorf("Expected documentId 'home', got '%s'", resp.Page.DocumentID)
	}
	
	if resp.Page.URL != "/" {
		t.Errorf("Expected URL '/', got '%s'", resp.Page.URL)
	}
	
	if len(resp.Page.Content) != 2 {
		t.Errorf("Expected 2 content components, got %d", len(resp.Page.Content))
	}
	
	if resp.Page.Metadata.Title != "Home Page" {
		t.Errorf("Expected metadata title 'Home Page', got '%s'", resp.Page.Metadata.Title)
	}
}

func TestPageNotFound(t *testing.T) {
	c := setupTestServer(t)
	
	var resp struct {
		Page interface{} `json:"page"`
	}
	
	query := `
		query {
			page(documentId: "non-existent") {
				documentId
			}
		}
	`
	
	err := c.Post(query, &resp)
	if err == nil {
		t.Error("Expected error for non-existent page, got none")
	}
}

func TestMenuItemsQuery(t *testing.T) {
	c := setupTestServer(t)
	
	var resp struct {
		MenuItems []struct {
			DocumentID string `json:"documentId"`
			Label      string `json:"label"`
			To         string `json:"to"`
			Items      []struct {
				DocumentID string `json:"documentId"`
				Label      string `json:"label"`
				To         string `json:"to"`
			} `json:"items"`
		} `json:"menuItems"`
	}
	
	query := `
		query {
			menuItems(pagination: {limit: 1000}) {
				documentId
				label
				to
				items {
					documentId
					label
					to
				}
			}
		}
	`
	
	err := c.Post(query, &resp)
	if err != nil {
		t.Fatalf("MenuItems query failed: %v", err)
	}
	
	if len(resp.MenuItems) != 1 {
		t.Errorf("Expected 1 menu root, got %d", len(resp.MenuItems))
	}
	
	if resp.MenuItems[0].Label != "$MAIN" {
		t.Errorf("Expected root label '$MAIN', got '%s'", resp.MenuItems[0].Label)
	}
	
	if len(resp.MenuItems[0].Items) != 2 {
		t.Errorf("Expected 2 menu items, got %d", len(resp.MenuItems[0].Items))
	}
}

func TestMenuItemsPagination(t *testing.T) {
	c := setupTestServer(t)
	
	var resp struct {
		MenuItems []interface{} `json:"menuItems"`
	}
	
	query := `
		query {
			menuItems(pagination: {limit: 0, offset: 10}) {
				documentId
			}
		}
	`
	
	err := c.Post(query, &resp)
	if err != nil {
		t.Fatalf("MenuItems pagination query failed: %v", err)
	}
	
	if len(resp.MenuItems) != 0 {
		t.Errorf("Expected 0 menu items with offset beyond count, got %d", len(resp.MenuItems))
	}
}

func TestFooterQuery(t *testing.T) {
	c := setupTestServer(t)
	
	var resp struct {
		Footer struct {
			ShortDescriptionTitle string `json:"shortDescriptionTitle"`
			ShortDescriptionText  string `json:"shortDescriptionText"`
			Links1                []struct {
				Label string `json:"label"`
				URL   string `json:"url"`
			} `json:"links_1"`
			Links2 []struct {
				Label string `json:"label"`
				URL   string `json:"url"`
			} `json:"links_2"`
			SocialLinks []struct {
				Label string `json:"label"`
				URL   string `json:"url"`
			} `json:"social_links"`
		} `json:"footer"`
	}
	
	query := `
		query {
			footer {
				shortDescriptionTitle
				shortDescriptionText
				links_1 {
					label
					url
				}
				links_2 {
					label
					url
				}
				social_links {
					label
					url
				}
			}
		}
	`
	
	err := c.Post(query, &resp)
	if err != nil {
		t.Fatalf("Footer query failed: %v", err)
	}
	
	if resp.Footer.ShortDescriptionTitle != "Test Site" {
		t.Errorf("Expected footer title 'Test Site', got '%s'", resp.Footer.ShortDescriptionTitle)
	}
	
	if len(resp.Footer.Links1) != 1 {
		t.Errorf("Expected 1 link in links_1, got %d", len(resp.Footer.Links1))
	}
	
	if len(resp.Footer.SocialLinks) != 1 {
		t.Errorf("Expected 1 social link, got %d", len(resp.Footer.SocialLinks))
	}
}

func TestComponentUnionTypes(t *testing.T) {
	c := setupTestServer(t)
	
	// Test that union types work correctly with __typename
	var resp map[string]interface{}
	
	query := `
		query {
			page(documentId: "home") {
				content {
					__typename
					... on ComponentPageHero {
						id
						title
					}
					... on ComponentPageMarkdown {
						id
						markdown
					}
				}
			}
		}
	`
	
	err := c.Post(query, &resp)
	if err != nil {
		t.Fatalf("Component union query failed: %v", err)
	}
	
	page := resp["page"].(map[string]interface{})
	content := page["content"].([]interface{})
	
	if len(content) != 2 {
		t.Errorf("Expected 2 content components, got %d", len(content))
	}
	
	// Check first component is Hero
	firstComponent := content[0].(map[string]interface{})
	if firstComponent["__typename"] != "ComponentPageHero" {
		t.Errorf("Expected first component type 'ComponentPageHero', got '%s'", firstComponent["__typename"])
	}
	
	// Check second component is Markdown
	secondComponent := content[1].(map[string]interface{})
	if secondComponent["__typename"] != "ComponentPageMarkdown" {
		t.Errorf("Expected second component type 'ComponentPageMarkdown', got '%s'", secondComponent["__typename"])
	}
}