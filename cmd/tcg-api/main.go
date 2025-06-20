package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
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

	// Set up structured logging
	logger := config.NewLogger(cfg.Logger)
	config.SetDefaultLogger(logger)

	logger.Info("Starting TCG API",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.Port))

	// Create a new HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      setupRoutes(*cfg, logger),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to start",
				slog.String("addr", server.Addr),
				slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Server exited")
}

func setupRoutes(cfg config.Config, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	// TODO: Initialize storage
	sto := storage.NewMockStorage()
	gameCardsHandler := handlers.NewGameCardsHandler(sto, logger)
	imageCardsHandler := handlers.NewImageCardsHandler(sto, logger)

	// Health endpoint
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Cards endpoints
	mux.Handle("/game-cards", gameCardsHandler)
	mux.Handle("/game-cards/", gameCardsHandler)

	mux.Handle("/image-cards", imageCardsHandler)
	mux.Handle("/image-cards/", imageCardsHandler)

	return mux
}
