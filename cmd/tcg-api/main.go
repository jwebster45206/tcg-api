package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jwebster45206/tcg-api/internal/config"
	"github.com/jwebster45206/tcg-api/internal/handlers"
	"github.com/jwebster45206/tcg-api/internal/storage"
)

// loadConfig loads configuration from config.json file
func loadConfig() (*config.Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Error closing config file: %v", closeErr)
		}
	}()

	var cfg config.Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a new HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      setupRoutes(*cfg),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRoutes(cfg config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags)

	// TODO: Initialize storage
	sto := storage.NewMockStorage()
	cardsHandler := handlers.NewGameCardsHandler(sto, logger)

	// Health endpoint
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Cards endpoints
	mux.Handle("/cards", cardsHandler)
	mux.Handle("/cards/", cardsHandler)

	return mux
}
