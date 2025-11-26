package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed all:web/dist
var webFS embed.FS

func main() {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// Serve frontend (embedded in prod, filesystem in dev)
	if err := serveFrontend(mux); err != nil {
		log.Printf("Warning: frontend not available: %v", err)
	}

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func serveFrontend(mux *http.ServeMux) error {
	var frontend fs.FS
	var err error

	// Try embedded first (production)
	frontend, err = fs.Sub(webFS, "web/dist")

	if err != nil {
		// Fallback to filesystem (development)
		if _, statErr := os.Stat("web/dist"); statErr == nil {
			frontend = os.DirFS("web/dist")
		} else {
			return err
		}
	}

	// Check if embedded fs has content (empty in dev)
	if entries, _ := fs.ReadDir(frontend, "."); len(entries) == 0 {
		// Fallback to filesystem (development)
		if _, statErr := os.Stat("web/dist"); statErr == nil {
			frontend = os.DirFS("web/dist")
		}
	}

	// Serve static files with SPA fallback
	fileServer := http.FileServer(http.FS(frontend))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists
		if f, err := frontend.Open(path[1:]); err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)

			return
		}

		// SPA fallback - serve index.html
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	return nil
}
