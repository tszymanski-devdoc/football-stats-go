# Football Analysis API

A lightweight Go API for analyzing football match data and predicting outcomes.

## Features

- **Team Statistics Analysis** - Calculate win/loss records, goals, and performance metrics
- **Match Prediction** - Predict match outcomes based on team statistics
- **Head-to-Head Analysis** - Analyze historical matchups between teams
- **Lightweight & Fast** - No database required, works with in-memory data

## API Endpoints

### 1. Analyze Team Statistics
```http
POST /api/analyze-team
Content-Type: application/json

{
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
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "team_name": "Manchester United",
    "matches_played": 10,
    "wins": 6,
    "draws": 2,
    "losses": 2,
    "goals_for": 18,
    "goals_against": 12,
    "win_percentage": 60.0,
    "avg_goals_scored": 1.8,
    "avg_goals_conceded": 1.2
  }
}
```

### 2. Predict Match Outcome
```http
POST /api/predict-match
Content-Type: application/json

{
  "home_team": "Manchester United",
  "away_team": "Chelsea",
  "home_matches": [...],
  "away_matches": [...]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "home_team": "Manchester United",
    "away_team": "Chelsea",
    "home_win_probability": 45.5,
    "draw_probability": 25.0,
    "away_win_probability": 29.5,
    "predicted_score": "2-1",
    "confidence": 75.0
  }
}
```

### 3. Head-to-Head Analysis
```http
POST /api/head-to-head
Content-Type: application/json

{
  "team1": "Manchester United",
  "team2": "Liverpool",
  "matches": [...]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "team1": "Manchester United",
    "team2": "Liverpool",
    "total_matches": 15,
    "team1_wins": 6,
    "team2_wins": 5,
    "draws": 4,
    "team1_avg_goals": 1.6,
    "team2_avg_goals": 1.4
  }
}
```

### 4. Health Check
```http
GET /health
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "football-analysis-api"
  }
}
```

## Quick Start

### Prerequisites
- Go 1.24+

### Run Locally

1. **Start the server:**
   ```bash
   go run cmd/api/main.go
   ```

2. **Test with curl:**
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Analyze team
   curl -X POST http://localhost:8080/api/analyze-team \
     -H "Content-Type: application/json" \
     -d '{
       "team_name": "Arsenal",
       "matches": [
         {
           "id": "1",
           "home_team": "Arsenal",
           "away_team": "Tottenham",
           "home_score": 3,
           "away_score": 1,
           "match_date": "2024-01-15T15:00:00Z",
           "league": "Premier League",
           "season": "2024"
         }
       ]
     }'
   ```

## Project Structure

```
football-stats-go-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/
│   │   └── models.go            # Data models
│   ├── analysis/
│   │   └── service.go           # Analysis logic
│   └── api/
│       ├── handler.go           # HTTP handlers
│       └── middleware.go        # HTTP middleware
└── internal/config/
    └── config.go                # Configuration
```

## Configuration

Configure via environment variables:

```bash
export APP_NAME=football-analytics
export APP_ENV=development
export SERVER_PORT=8080
```

## Example Usage

### Python Example
```python
import requests

# Analyze team
response = requests.post('http://localhost:8080/api/analyze-team', json={
    'team_name': 'Barcelona',
    'matches': [
        {
            'id': '1',
            'home_team': 'Barcelona',
            'away_team': 'Real Madrid',
            'home_score': 2,
            'away_score': 1,
            'match_date': '2024-01-15T15:00:00Z',
            'league': 'La Liga',
            'season': '2024'
        }
    ]
})

print(response.json())
```

### JavaScript Example
```javascript
fetch('http://localhost:8080/api/predict-match', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    home_team: 'Liverpool',
    away_team: 'Manchester City',
    home_matches: [...],
    away_matches: [...]
  })
})
.then(res => res.json())
.then(data => console.log(data));
```

## Deployment to GCP

### Cloud Run
```bash
# Build and deploy
gcloud run deploy football-analysis-api \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```

### Using Dockerfile
```bash
# Build
docker build -t football-analysis-api .

# Run locally
docker run -p 8080:8080 football-analysis-api

# Deploy to GCP
gcloud builds submit --tag gcr.io/PROJECT_ID/football-analysis-api
gcloud run deploy --image gcr.io/PROJECT_ID/football-analysis-api
```

## Features

✅ Lightweight - No database required  
✅ RESTful API - Clean JSON endpoints  
✅ CORS enabled - Works with web frontends  
✅ Request logging - Track all API calls  
✅ Graceful shutdown - Proper cleanup on exit  
✅ Health checks - Monitor service status  
✅ Cloud-ready - Deploy to GCP Cloud Run  

## License

MIT
