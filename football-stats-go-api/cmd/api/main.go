package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpSwagger "github.com/swaggo/http-swagger"

	"example/hello/internal/analysis"
	"example/hello/internal/api"
	"example/hello/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("Starting %s v%s in %s environment", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Initialize analysis service
	analysisService := analysis.NewService()

	// Setup API handler
	apiHandler := api.NewHandler(analysisService)

	// Setup HTTP router
	mux := http.NewServeMux()
	mux.HandleFunc("/health", apiHandler.Health)
	mux.HandleFunc("/api/analyze-team", apiHandler.AnalyzeTeam)
	mux.HandleFunc("/api/predict-match", apiHandler.PredictMatch)
	mux.HandleFunc("/api/head-to-head", apiHandler.HeadToHead)

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Apply middleware
	handler := api.LoggingMiddleware(api.CORSMiddleware(mux))

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		log.Println("Available endpoints:")
		log.Println("  POST /api/analyze-team     - Analyze team statistics")
		log.Printf("  GET  /swagger/             - Swagger UI (http://%s/swagger/)\n", addr)
		log.Println("  POST /api/predict-match    - Predict match outcome")
		log.Println("  POST /api/head-to-head     - Head-to-head analysis")
		log.Println("  GET  /health               - Health check")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
