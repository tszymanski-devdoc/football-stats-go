package domain

import "time"

type DBXGStatShot struct {
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	XG         float64 `json:"xg"`
	IsGoal     bool    `json:"is_goal"`
	ShotType   string  `json:"shot_type"`
	PlayerName string  `json:"player_name"`
	Minute     int     `json:"minute"`
}

type DBXGStatFixture struct {
	Gameweek  int            `json:"gameweek"`
	ID        int            `json:"id"`
	Date      time.Time      `json:"date"`
	HomeTeam  string         `json:"home_team"`
	AwayTeam  string         `json:"away_team"`
	HomeScore int            `json:"home_score"`
	AwayScore int            `json:"away_score"`
	HomeXG    float64        `json:"home_xg"`
	AwayXG    float64        `json:"away_xg"`
	HomeShots []DBXGStatShot `json:"home_shots"`
	AwayShots []DBXGStatShot `json:"away_shots"`
}
