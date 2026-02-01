# xG Shot Map Scraping Guide

This guide explains how to scrape Expected Goals (xG) shot map data from xgstat.com.

## Overview

The API now supports scraping detailed shot map data including:
- Match details (teams, score, date, gameweek)
- xG values for both teams
- Individual shot data with coordinates, xG values, and outcomes
- Player names and shot types

## Data Model

### DBXGStatFixture
```go
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
```

### DBXGStatShot
```go
type DBXGStatShot struct {
    X          float64 `json:"x"`          // X coordinate on pitch (0-100)
    Y          float64 `json:"y"`          // Y coordinate on pitch (0-100)
    XG         float64 `json:"xg"`         // Expected goals value (0-1)
    IsGoal     bool    `json:"is_goal"`    // Whether shot resulted in goal
    ShotType   string  `json:"shot_type"`  // e.g., "Right Foot", "Header"
    PlayerName string  `json:"player_name"`// Name of the player
    Minute     int     `json:"minute"`     // Minute of the match
}
```

## Usage

### API Endpoint
```
POST /api/scrape/xgstats
```

### Request Format
```json
{
  "url": "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"
}
```

### Response Format
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
    "away_shots": [...]
  }
}
```

## Finding Match URLs

xgstat.com URLs follow this pattern:
```
https://www.xgstat.com/competitions/{competition}/{season}/matches/{match-slug}/advanced-analysis/shot-maps
```

Examples:
- Premier League: `premier-league/2025-2026`
- La Liga: `la-liga/2025-2026`
- Champions League: `champions-league/2025-2026`

## Testing

Use the provided test scripts:

### PowerShell
```powershell
.\test-xgstats.ps1
```

### Bash
```bash
chmod +x test-xgstats.sh
./test-xgstats.sh
```

### cURL
```bash
curl -X POST http://localhost:8080/api/scrape/xgstats \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"
  }'
```

## Environment Variables

Configure the scraper behavior:

```bash
# Run browser in visible mode (for debugging)
export SCRAPER_HEADLESS=false

# Enable debug logging
export SCRAPER_DEBUG=true
```

## Implementation Details

### Scraping Process

1. **Page Navigation**: Uses headless Chrome to navigate to the match URL
2. **Data Extraction**: Extracts the `window.__NUXT__` JavaScript object containing all match data
3. **Parsing**: Parses the JSON data to extract match details and shot coordinates
4. **Validation**: Validates and structures data into `DBXGStatFixture` model

### Key Features

- **Bot Detection Bypass**: Uses realistic browser headers and behavior
- **Timeout Handling**: 60-second timeout for page loading
- **Error Handling**: Comprehensive error messages for debugging
- **Flexible Parsing**: Handles various data structures from xgstat.com

## Common Issues

### Page Not Loading
- Increase timeout in [service.go](service.go#L60)
- Check if URL is correct and match exists
- Enable visible browser mode: `SCRAPER_HEADLESS=false`

### Missing Data
- Some matches may not have complete shot data
- Check the actual webpage to verify data availability
- Enable debug mode: `SCRAPER_DEBUG=true`

### Bot Detection
- xgstat.com may update their bot detection
- The scraper includes anti-bot measures, but these may need updates
- Try running in non-headless mode first

## Example Use Cases

### 1. Match Analysis
Analyze shot quality and positioning for tactical insights

### 2. Player Performance
Track individual player xG contribution across matches

### 3. Team Statistics
Aggregate shot data to evaluate team offensive patterns

### 4. Historical Data
Build a database of historical shot map data for analysis

## Next Steps

Consider adding:
- Database storage for scraped fixtures
- Batch scraping for multiple matches
- API endpoints to query stored data
- Data visualization endpoints
- Scheduled scraping jobs
