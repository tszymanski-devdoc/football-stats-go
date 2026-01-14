package api

import (
	"encoding/json"
	"net/http"

	"example/hello/internal/analysis"
	"example/hello/internal/domain"
)

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
