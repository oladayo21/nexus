package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"

	nexus "github.com/oladayo21/nexus"
)

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
	frontend, err = nexus.WebFS()

	if err != nil {
		// Fallback to filesystem (development)
		if _, statErr := os.Stat("web/dist"); statErr == nil {
			frontend = os.DirFS("web/dist")
		} else {
			return err
		}
	}

	// Serve static files with SPA fallback
	fileServer := http.FileServer(http.FS(frontend))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file
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

		// SPA fallback - serve index.html for non-API routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	return nil
}
