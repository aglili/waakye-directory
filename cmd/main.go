package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aglili/waakye-directory/internal/config"
	"github.com/aglili/waakye-directory/internal/logger"
	"github.com/aglili/waakye-directory/internal/provider"
	"github.com/aglili/waakye-directory/internal/routes"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	logger.Init(config.GetEnvOrDefault("ENV", "development"))
	log.Info().Msg("Starting the application...")

	// Load configuration
	cfg := config.LoadConfig()
	log.Info().Interface("config", cfg).Msg("Loaded configuration")

	// Initialize database
	db, err := config.InitializeDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize the database")
	}
	defer db.Close()

	// Initialize cache
	cache, err := config.InitializeCache(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize the cache")
	}

	// Create a new provider
	prov := provider.NewProvider(db, cache)

	// Setup routes
	router := routes.SetupRoutes(prov)

	// Get server port
	serverPort := fmt.Sprintf(":%s", config.GetEnvOrDefault("PORT", "8080"))

	// Create HTTP server
	httpServer := &http.Server{
		Addr:           serverPort,
		Handler:        router,
		MaxHeaderBytes: 1 << 20,           // 1 MB
		ReadTimeout:    10 * time.Second,  // 10 seconds
		WriteTimeout:   10 * time.Second,  // 10 seconds
		IdleTimeout:    120 * time.Second, // 2 minutes
	}

	// Start server in a separate goroutine
	go func() {
		log.Info().Msgf("Starting server on port %s", serverPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	// Graceful shutdown
	gracefulShutdown(httpServer)
}

func gracefulShutdown(server *http.Server) {
	// Create channel to receive OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-stop
	log.Info().Msg("Shutting down server gracefully...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server stopped gracefully")
}
