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
	"github.com/redis/go-redis/v9"
)

//go:embed all:web/dist
var webFS embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	staticFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to get static filesystem: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	apiServer := api.NewAPIServer(&api.Options{
		Config: cfg,
		Redis:  redisClient,
		WebFS:  staticFS,
	})

	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	log.Printf("Server started on :%d", cfg.Port)

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
