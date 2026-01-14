package analysis

import (
	"fmt"
	"math"

	"example/hello/internal/domain"
)

// Service provides football data analysis capabilities
type Service struct {
	// In a real app, this would have database repos or data sources
}

// NewService creates a new analysis service
func NewService() *Service {
	return &Service{}
}

// CalculateTeamStats computes statistics for a team based on match data
func (s *Service) CalculateTeamStats(teamName string, matches []domain.MatchData) *domain.TeamStats {
	stats := &domain.TeamStats{
		TeamName: teamName,
	}

	for _, match := range matches {
		isHome := match.HomeTeam == teamName
		isAway := match.AwayTeam == teamName

		if !isHome && !isAway {
			continue
		}

		stats.MatchesPlayed++

		var goalsFor, goalsAgainst int
		if isHome {
			goalsFor = match.HomeScore
			goalsAgainst = match.AwayScore
		} else {
			goalsFor = match.AwayScore
			goalsAgainst = match.HomeScore
		}

		stats.GoalsFor += goalsFor
		stats.GoalsAgainst += goalsAgainst

		if goalsFor > goalsAgainst {
			stats.Wins++
		} else if goalsFor == goalsAgainst {
			stats.Draws++
		} else {
			stats.Losses++
		}
	}

	if stats.MatchesPlayed > 0 {
		stats.WinPercentage = float64(stats.Wins) / float64(stats.MatchesPlayed) * 100
		stats.AvgGoalsScored = float64(stats.GoalsFor) / float64(stats.MatchesPlayed)
		stats.AvgGoalsConceded = float64(stats.GoalsAgainst) / float64(stats.MatchesPlayed)
	}

	return stats
}

// PredictMatch generates a simple prediction based on team statistics
func (s *Service) PredictMatch(homeStats, awayStats *domain.TeamStats) *domain.MatchPrediction {
	// Simple prediction model based on win percentages and average goals
	homeStrength := homeStats.WinPercentage + (homeStats.AvgGoalsScored * 10) - (homeStats.AvgGoalsConceded * 5)
	awayStrength := awayStats.WinPercentage + (awayStats.AvgGoalsScored * 10) - (awayStats.AvgGoalsConceded * 5)

	// Home advantage factor
	homeStrength += 10

	total := homeStrength + awayStrength
	if total == 0 {
		total = 1
	}

	homeWinProb := (homeStrength / total) * 100
	awayWinProb := (awayStrength / total) * 100
	drawProb := 100 - homeWinProb - awayWinProb

	// Ensure probabilities are reasonable
	if drawProb < 0 {
		drawProb = 15
		homeWinProb = (homeWinProb / (homeWinProb + awayWinProb)) * 85
		awayWinProb = 85 - homeWinProb
	}

	// Predict score based on average goals
	homeGoals := math.Round(homeStats.AvgGoalsScored)
	awayGoals := math.Round(awayStats.AvgGoalsScored * 0.8) // Reduce away team goals slightly

	predictedScore := fmt.Sprintf("%.0f-%.0f", homeGoals, awayGoals)

	// Calculate confidence based on data quality
	confidence := calculateConfidence(homeStats, awayStats)

	return &domain.MatchPrediction{
		HomeTeam:           homeStats.TeamName,
		AwayTeam:           awayStats.TeamName,
		HomeWinProbability: math.Round(homeWinProb*100) / 100,
		DrawProbability:    math.Round(drawProb*100) / 100,
		AwayWinProbability: math.Round(awayWinProb*100) / 100,
		PredictedScore:     predictedScore,
		Confidence:         confidence,
	}
}

// AnalyzeHeadToHead analyzes historical matchups between two teams
func (s *Service) AnalyzeHeadToHead(team1, team2 string, matches []domain.MatchData) *domain.HeadToHeadStats {
	h2h := &domain.HeadToHeadStats{
		Team1: team1,
		Team2: team2,
	}

	var team1Goals, team2Goals int

	for _, match := range matches {
		// Check if both teams are in this match
		if (match.HomeTeam == team1 && match.AwayTeam == team2) ||
			(match.HomeTeam == team2 && match.AwayTeam == team1) {

			h2h.TotalMatches++

			var t1Score, t2Score int
			if match.HomeTeam == team1 {
				t1Score = match.HomeScore
				t2Score = match.AwayScore
			} else {
				t1Score = match.AwayScore
				t2Score = match.HomeScore
			}

			team1Goals += t1Score
			team2Goals += t2Score

			if t1Score > t2Score {
				h2h.Team1Wins++
			} else if t2Score > t1Score {
				h2h.Team2Wins++
			} else {
				h2h.Draws++
			}
		}
	}

	if h2h.TotalMatches > 0 {
		h2h.Team1AvgGoals = float64(team1Goals) / float64(h2h.TotalMatches)
		h2h.Team2AvgGoals = float64(team2Goals) / float64(h2h.TotalMatches)
	}

	return h2h
}

// calculateConfidence determines prediction confidence based on data quality
func calculateConfidence(homeStats, awayStats *domain.TeamStats) float64 {
	confidence := 50.0

	// More matches = higher confidence
	avgMatches := float64(homeStats.MatchesPlayed+awayStats.MatchesPlayed) / 2
	if avgMatches >= 10 {
		confidence += 20
	} else if avgMatches >= 5 {
		confidence += 10
	}

	// Consistent performance = higher confidence
	if homeStats.WinPercentage > 60 || awayStats.WinPercentage > 60 {
		confidence += 15
	}

	return math.Min(confidence, 95.0)
}
