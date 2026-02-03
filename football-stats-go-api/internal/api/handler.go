package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example/hello/internal/database"
	"example/hello/internal/scraper"
)

// Handler handles HTTP requests for the scraper API
type Handler struct {
	scraperService  *scraper.Service
	databaseService *database.Service
}

// NewHandler creates a new API handler
func NewHandler(scraperService *scraper.Service, databaseService *database.Service) *Handler {
	return &Handler{
		scraperService:  scraperService,
		databaseService: databaseService,
	}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ScrapeRequest represents a request to scrape a website
type ScrapeRequest struct {
	URL string `json:"url"`
}

// ScrapeXGStats scrapes xG shot map data from xgstat.com
// @Summary Scrape xG shot map data
// @Description Scrape xG statistics and shot map data from xgstat.com
// @Tags scraper
// @Accept json
// @Produce json
// @Param request body ScrapeRequest true "Scrape request with xgstat.com URL"
// @Success 200 {object} Response{data=example_hello_internal_domain.DBXGStatFixture} "Scraped xG statistics and shot map data"
// @Failure 400 {object} Response "Invalid request"
// @Failure 405 {object} Response "Method not allowed"
// @Router /scrape/xgstats [post]
func (h *Handler) ScrapeXGStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "URL is required")
		return
	}

	data, err := h.scraperService.ScrapeXGStatFixture(req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Save to database if service is available
	if h.databaseService != nil {
		if err := h.databaseService.SaveXGStatFixture(data); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to save data: "+err.Error())
			return
		}
	}

	writeSuccess(w, data)
}

// GetXGStatFixture retrieves a saved fixture by ID
// @Summary Get xG statistics by fixture ID
// @Description Retrieve saved xG statistics and shot map data from the database
// @Tags scraper
// @Accept json
// @Produce json
// @Param id query int true "Fixture ID"
// @Success 200 {object} Response{data=example_hello_internal_domain.DBXGStatFixture} "Retrieved xG statistics"
// @Failure 400 {object} Response "Invalid request"
// @Failure 404 {object} Response "Fixture not found"
// @Router /xgstats [get]
func (h *Handler) GetXGStatFixture(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	fixtureIDStr := r.URL.Query().Get("id")
	if fixtureIDStr == "" {
		writeError(w, http.StatusBadRequest, "Fixture ID is required")
		return
	}

	var fixtureID int
	if _, err := fmt.Sscanf(fixtureIDStr, "%d", &fixtureID); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid fixture ID")
		return
	}

	data, err := h.databaseService.GetFixtureByID(fixtureID)
	if err != nil {
		if err.Error() == "fixture not found" {
			writeError(w, http.StatusNotFound, "Fixture not found")
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeSuccess(w, data)
}

// Health returns the health status
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeSuccess(w, map[string]string{
		"status":  "healthy",
		"service": "football-scraper-api",
	})
}

// writeSuccess writes a successful JSON response
func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

// writeError writes an error JSON response
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   message,
	})
}
