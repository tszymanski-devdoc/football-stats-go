package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"example/hello/internal/api"
	"example/hello/internal/config"
	"example/hello/internal/database"
	"example/hello/internal/scraper"

	_ "example/hello/docs"
)

// @title Football Stats Scraper API
// @version 1.0
// @description A lightweight API for scraping football data from websites
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@footballstats.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

func main() {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	log.Printf("Starting %s v%s in %s environment", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Log important environment info for Cloud Run debugging
	log.Printf("PORT: %s", cfg.Server.Port)
	log.Printf("DATABASE_URL configured: %v", os.Getenv("DATABASE_URL") != "")
	log.Printf("CHROME_PATH: %s", os.Getenv("CHROME_PATH"))

	// Initialize database service
	dbService, err := database.NewService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Continuing without database connection...")
		dbService = nil
	} else {
		log.Println("Database connection established")
	}

	// Initialize scraper service
	scraperService := scraper.NewService()

	// Setup API handler
	apiHandler := api.NewHandler(scraperService, dbService)

	// Setup HTTP router
	mux := http.NewServeMux()
	mux.HandleFunc("/health", apiHandler.Health)
	mux.HandleFunc("/api/scrape/xgstats", apiHandler.ScrapeXGStats)
	mux.HandleFunc("/api/xgstats", apiHandler.GetXGStatFixture)

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
		log.Println("  POST /api/scrape/xgstats   - Scrape xG shot map data from xgstat.com")
		log.Println("  GET  /api/xgstats?id=XXX   - Get saved xG statistics by fixture ID")
		log.Println("  GET  /health               - Health check")
		log.Printf("  GET  /swagger/             - Swagger UI (http://%s/swagger/)\n", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if dbService != nil {
		dbService.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
