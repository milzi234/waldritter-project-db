package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/waldritter/website-content-api/content"
	"github.com/waldritter/website-content-api/graph"
	"github.com/waldritter/website-content-api/graph/generated"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "1337", "port to run the server on")
	vaultPath := flag.String("vault", "./obsidian-vault", "path to the Obsidian vault")
	enablePlayground := flag.Bool("playground", true, "enable GraphQL playground")
	enableWatcher := flag.Bool("watch", true, "enable file watcher for hot reload")
	flag.Parse()

	// Initialize content loader and store
	loader := content.NewContentLoader(*vaultPath)
	store := content.NewContentStore()

	// Initial content load
	log.Println("Loading content from vault...")
	data, err := loader.LoadContent()
	if err != nil {
		log.Printf("Warning: Failed to load initial content: %v", err)
	} else {
		store.UpdateContent(data.Pages, data.MenuItems, data.Footer)
		stats := store.GetStats()
		log.Printf("Content loaded successfully: %d pages, %d menu items, footer=%v",
			stats["pages_count"], stats["menu_items"], stats["has_footer"])
	}

	// Validate content
	validationErrors := loader.ValidateContent(data)
	if len(validationErrors) > 0 {
		log.Println("Content validation warnings:")
		for _, err := range validationErrors {
			log.Printf("  - %s", err)
		}
	}

	// Start file watcher if enabled
	if *enableWatcher {
		watcher, err := content.NewContentWatcher(*vaultPath, loader, store)
		if err != nil {
			log.Printf("Warning: Failed to start file watcher: %v", err)
		} else {
			err = watcher.Start()
			if err != nil {
				log.Printf("Warning: Failed to start watcher: %v", err)
			} else {
				log.Println("File watcher started - content will hot-reload on changes")
				defer watcher.Stop()
			}
		}
	}

	// Create GraphQL server
	resolver := &graph.Resolver{
		Store: store,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // Admin UI dev server
			"http://localhost:5174", // Public UI dev server
			"http://localhost:3000", // Rails API
			"http://localhost:8080", // Alternative port
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
			"Accept",
			"Accept-Language",
			"Accept-Encoding",
			"Content-Length",
			"Origin",
		},
		AllowCredentials: true,
		Debug:            false,
	})

	// Set up routes
	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", c.Handler(srv))

	// GraphQL playground (if enabled)
	if *enablePlayground {
		mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		log.Println("GraphQL playground enabled at http://localhost:" + *port + "/")
	}

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		stats := store.GetStats()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","pages":%d,"version":"%s"}`,
			stats["pages_count"], stats["version"])
	})

	// Static assets handler
	assetsPath := fmt.Sprintf("%s/assets", *vaultPath)
	if _, err := os.Stat(assetsPath); err == nil {
		fs := http.FileServer(http.Dir(assetsPath))
		mux.Handle("/assets/", http.StripPrefix("/assets/", c.Handler(fs)))
		log.Printf("Serving static assets from %s", assetsPath)
	}
	
	// Uploads handler for backward compatibility with Strapi URLs
	uploadsPath := fmt.Sprintf("%s/assets/images/uploads", *vaultPath)
	if _, err := os.Stat(uploadsPath); err == nil {
		fs := http.FileServer(http.Dir(uploadsPath))
		mux.Handle("/uploads/", http.StripPrefix("/uploads/", c.Handler(fs)))
		log.Printf("Serving uploads from %s", uploadsPath)
	}

	// Start server
	log.Printf("Starting GraphQL server on port %s", *port)
	log.Printf("GraphQL endpoint: http://localhost:%s/graphql", *port)
	
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}