package api

import (
	"encoding/json"
	"net/http"

	"example/hello/internal/analysis"
	"example/hello/internal/domain"
)

// @title Football Stats Analysis API
// @version 1.0
// @description A lightweight API for analyzing football data and predicting match outcomes
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@footballstats.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

// Handler handles HTTP requests for the analysis API
type Handler struct {
	analysisService *analysis.Service
}

// NewHandler creates a new API handler
func NewHandler(analysisService *analysis.Service) *Handler {
	return &Handler{
		analysisService: analysisService,
	}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// AnalyzeTeamRequest represents a request to analyze team statistics
type AnalyzeTeamRequest struct {
	TeamName string             `json:"team_name"`
	Matches  []domain.MatchData `json:"matches"`
}

// PredictMatchRequest represents a request to predict a match
type PredictMatchRequest struct {
	HomeTeam    string             `json:"home_team"`
	AwayTeam    string             `json:"away_team"`
	HomeMatches []domain.MatchData `json:"home_matches"`
	AwayMatches []domain.MatchData `json:"away_matches"`
}

// HeadToHeadRequest represents a request for head-to-head analysis
type HeadToHeadRequest struct {
	Team1   string             `json:"team1"`
	Team2   string             `json:"team2"`
	Matches []domain.MatchData `json:"matches"`
}

// AnalyzeTeam analyzes team statistics
// @Summary Analyze team statistics
// @Description Calculate comprehensive statistics for a team based on match data
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body AnalyzeTeamRequest true "Team analysis request"
// @Success 200 {object} Response{data=domain.TeamStats} "Team statistics"
// @Failure 400 {object} Response "Invalid request"
// @Failure 405 {object} Response "Method not allowed"
// @Router /analyze-team [post]
func (h *Handler) AnalyzeTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req AnalyzeTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TeamName == "" {
		writeError(w, http.StatusBadRequest, "team_name is required")
		return
	}

	stats := h.analysisService.CalculateTeamStats(req.TeamName, req.Matches)
	writeSuccess(w, stats)
}

// PredictMatch predicts match outcome
// @Summary Predict match outcome
// @Description Predict the outcome of a match based on team statistics using Poisson distribution
// @Tags prediction
// @Accept json
// @Produce json
// @Param request body PredictMatchRequest true "Match prediction request"
// @Success 200 {object} Response{data=domain.MatchPrediction} "Match prediction"
// @Failure 400 {object} Response "Invalid request"
// @Failure 405 {object} Response "Method not allowed"
// @Router /predict-match [post]
func (h *Handler) PredictMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req PredictMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.HomeTeam == "" || req.AwayTeam == "" {
		writeError(w, http.StatusBadRequest, "home_team and away_team are required")
		return
		// @Summary Head-to-head analysis
		// @Description Analyze historical matchups between two teams
		// @Tags analysis
		// @Accept json
		// @Produce json
		// @Param request body HeadToHeadRequest true "Head-to-head request"
		// @Success 200 {object} Response{data=domain.HeadToHeadStats} "Head-to-head statistics"
		// @Failure 400 {object} Response "Invalid request"
		// @Failure 405 {object} Response "Method not allowed"
		// @Router /head-to-head [post]
	}

	homeStats := h.analysisService.CalculateTeamStats(req.HomeTeam, req.HomeMatches)
	awayStats := h.analysisService.CalculateTeamStats(req.AwayTeam, req.AwayMatches)

	prediction := h.analysisService.PredictMatch(homeStats, awayStats)
	writeSuccess(w, prediction)
}

// HeadToHead analyzes head-to-head statistics
func (h *Handler) HeadToHead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req HeadToHeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// @Summary Health check
	// @Description Returns the health status of the API
	// @Tags system
	// @Produce json
	// @Success 200 {object} Response{data=map[string]string} "Health status"
	// @Router /health [get]
	if req.Team1 == "" || req.Team2 == "" {
		writeError(w, http.StatusBadRequest, "team1 and team2 are required")
		return
	}

	stats := h.analysisService.AnalyzeHeadToHead(req.Team1, req.Team2, req.Matches)
	writeSuccess(w, stats)
}

// Health returns the health status
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeSuccess(w, map[string]string{
		"status":  "healthy",
		"service": "football-analysis-api",
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
