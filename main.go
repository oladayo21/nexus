package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oladayo21/nexus/internal/api"
	"github.com/oladayo21/nexus/internal/config"
)

//go:embed all:web/dist
var webFS embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get sub filesystem for web/dist
	staticFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to get static filesystem: %v", err)
	}

	apiServer := api.NewAPIServer(&api.APIServerOptions{
		Port:        cfg.Port,
		ServeStatic: cfg.IsProduction(),
		WebFS:       staticFS,
	})

	// Start server in goroutine
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	log.Printf("Server started on :%d", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Stop(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
