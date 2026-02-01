# Football Scraper API

A lightweight Go API for scraping football data from websites.

## Features

- **Web Scraping** - Scrape data from football websites using headless Chrome
- **Lightweight & Fast** - Simple and efficient scraping service
- **REST API** - Easy to use HTTP endpoints

## API Endpoints

### 1. Scrape xG Shot Map Data
```http
POST /api/scrape/xgstats
Content-Type: application/json

{
  "url": "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "gameweek": 23,
    "id": 12345,
    "date": "2026-01-24T15:00:00Z",
    "home_team": "Arsenal",
    "away_team": "Manchester United",
    "home_score": 2,
    "away_score": 1,
    "home_xg": 2.34,
    "away_xg": 1.12,
    "home_shots": [
      {
        "x": 88.5,
        "y": 45.2,
        "xg": 0.45,
        "is_goal": true,
        "shot_type": "Right Foot",
        "player_name": "Bukayo Saka",
        "minute": 23
      }
    ],
    "away_shots": [
      {
        "x": 12.3,
        "y": 50.1,
        "xg": 0.32,
        "is_goal": false,
        "shot_type": "Left Foot",
        "player_name": "Marcus Rashford",
        "minute": 67
      }
    ]
  }
}
```

### 2. Scrape Website (Legacy)
```http
POST /api/scrape
Content-Type: application/json

{
  "url": "https://www.premierleague.com"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "url": "https://www.premierleague.com",
    "title": "Premier League - Official Website",
    "content": "...",
    "timestamp": "2026-01-26T10:30:00Z",
    "metadata": {
      "scraped_at": "2026-01-26T10:30:00Z",
      "status": "success"
    }
  }
}
```

### 3. Predict Match Outcome
```http
POST /api/predict-match
Content-Type: application/json

{
  "home_team": "Manchester United",
  "away_team": "Chelsea",
  "home_matches": [...],
```

### 2. Health Check
```http
GET /health
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "football-scraper-api"
  }
}
```

## Quick Start

### Prerequisites
- Go 1.24+
- Chrome/Chromium (for headless scraping)

### Run Locally

1. **Start the server:**
   ```bash
   go run cmd/api/main.go
   ```

2. **Test with curl:**
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Scrape a website
   curl -X POST http://localhost:8080/api/scrape \
     -H "Content-Type: application/json" \
     -d '{"url": "https://www.premierleague.com"}'
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
│   ├── scraper/
│   │   └── service.go           # Scraper logic
│   └── api/
│       ├── handler.go           # HTTP handlers
│       └── middleware.go        # HTTP middleware
└── internal/config/
    └── config.go                # Configuration
```

## Configuration

Configure via environment variables:

```bash
export APP_NAME=football-scraper
export APP_ENV=development
export SERVER_PORT=8080
```

## Example Usage

### Python Example
```python
import requests

# Scrape website
response = requests.post('http://localhost:8080/api/scrape', json={
    'url': 'https://www.premierleague.com'
})

print(response.json())
```

### JavaScript Example
```javascript
fetch('http://localhost:8080/api/scrape', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    url: 'https://www.premierleague.com'
  })
})
.then(res => res.json())
.then(data => console.log(data));
```

## Deployment to GCP

### Cloud Run
```bash
# Build and deploy
gcloud run deploy football-scraper-api \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```

### Using Dockerfile
```bash
# Build
docker build -t football-scraper-api .

# Run locally
docker run -p 8080:8080 football-scraper-api

# Deploy to GCP
gcloud builds submit --tag gcr.io/PROJECT_ID/football-scraper-api
gcloud run deploy --image gcr.io/PROJECT_ID/football-scraper-api
```

## Features

✅ Web Scraping - Extract data from websites  
✅ Headless Chrome - Fast and reliable scraping  
✅ RESTful API - Clean JSON endpoints  
✅ CORS enabled - Works with web frontends  
✅ Request logging - Track all API calls  
✅ Graceful shutdown - Proper cleanup on exit  
✅ Health checks - Monitor service status  
✅ Cloud-ready - Deploy to GCP Cloud Run  

## License

MIT
