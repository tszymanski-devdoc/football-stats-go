package domain

import "time"

// MatchData represents basic match information for analysis
type MatchData struct {
	ID        string    `json:"id"`
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	HomeScore int       `json:"home_score"`
	AwayScore int       `json:"away_score"`
	MatchDate time.Time `json:"match_date"`
	League    string    `json:"league"`
	Season    string    `json:"season"`
}

// TeamStats represents aggregated team statistics
type TeamStats struct {
	TeamName         string  `json:"team_name"`
	MatchesPlayed    int     `json:"matches_played"`
	Wins             int     `json:"wins"`
	Draws            int     `json:"draws"`
	Losses           int     `json:"losses"`
	GoalsFor         int     `json:"goals_for"`
	GoalsAgainst     int     `json:"goals_against"`
	WinPercentage    float64 `json:"win_percentage"`
	AvgGoalsScored   float64 `json:"avg_goals_scored"`
	AvgGoalsConceded float64 `json:"avg_goals_conceded"`
}

// MatchPrediction represents predicted match outcome
type MatchPrediction struct {
	HomeTeam           string  `json:"home_team"`
	AwayTeam           string  `json:"away_team"`
	HomeWinProbability float64 `json:"home_win_probability"`
	DrawProbability    float64 `json:"draw_probability"`
	AwayWinProbability float64 `json:"away_win_probability"`
	PredictedScore     string  `json:"predicted_score"`
	Confidence         float64 `json:"confidence"`
}

// HeadToHeadStats represents historical matchup statistics
type HeadToHeadStats struct {
	Team1         string  `json:"team1"`
	Team2         string  `json:"team2"`
	TotalMatches  int     `json:"total_matches"`
	Team1Wins     int     `json:"team1_wins"`
	Team2Wins     int     `json:"team2_wins"`
	Draws         int     `json:"draws"`
	Team1AvgGoals float64 `json:"team1_avg_goals"`
	Team2AvgGoals float64 `json:"team2_avg_goals"`
}
