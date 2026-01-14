# Test the Football Analysis API

# Health check
curl http://localhost:8080/health

# Analyze team statistics
curl -X POST http://localhost:8080/api/analyze-team \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "Manchester United",
    "matches": [
      {
        "id": "1",
        "home_team": "Manchester United",
        "away_team": "Liverpool",
        "home_score": 2,
        "away_score": 1,
        "match_date": "2024-01-15T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      },
      {
        "id": "2",
        "home_team": "Chelsea",
        "away_team": "Manchester United",
        "home_score": 1,
        "away_score": 1,
        "match_date": "2024-01-22T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      },
      {
        "id": "3",
        "home_team": "Manchester United",
        "away_team": "Arsenal",
        "home_score": 3,
        "away_score": 0,
        "match_date": "2024-01-29T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      }
    ]
  }'

# Predict match outcome
curl -X POST http://localhost:8080/api/predict-match \
  -H "Content-Type: application/json" \
  -d '{
    "home_team": "Manchester United",
    "away_team": "Liverpool",
    "home_matches": [
      {
        "id": "1",
        "home_team": "Manchester United",
        "away_team": "Chelsea",
        "home_score": 2,
        "away_score": 1,
        "match_date": "2024-01-15T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      },
      {
        "id": "2",
        "home_team": "Manchester United",
        "away_team": "Arsenal",
        "home_score": 3,
        "away_score": 0,
        "match_date": "2024-01-22T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      }
    ],
    "away_matches": [
      {
        "id": "3",
        "home_team": "Liverpool",
        "away_team": "Chelsea",
        "home_score": 4,
        "away_score": 1,
        "match_date": "2024-01-15T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      },
      {
        "id": "4",
        "home_team": "Liverpool",
        "away_team": "Arsenal",
        "home_score": 2,
        "away_score": 2,
        "match_date": "2024-01-22T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      }
    ]
  }'

# Head-to-head analysis
curl -X POST http://localhost:8080/api/head-to-head \
  -H "Content-Type: application/json" \
  -d '{
    "team1": "Manchester United",
    "team2": "Liverpool",
    "matches": [
      {
        "id": "1",
        "home_team": "Manchester United",
        "away_team": "Liverpool",
        "home_score": 2,
        "away_score": 1,
        "match_date": "2023-09-15T15:00:00Z",
        "league": "Premier League",
        "season": "2023"
      },
      {
        "id": "2",
        "home_team": "Liverpool",
        "away_team": "Manchester United",
        "home_score": 3,
        "away_score": 1,
        "match_date": "2023-12-10T15:00:00Z",
        "league": "Premier League",
        "season": "2023"
      },
      {
        "id": "3",
        "home_team": "Manchester United",
        "away_team": "Liverpool",
        "home_score": 1,
        "away_score": 1,
        "match_date": "2024-03-20T15:00:00Z",
        "league": "Premier League",
        "season": "2024"
      }
    ]
  }'
