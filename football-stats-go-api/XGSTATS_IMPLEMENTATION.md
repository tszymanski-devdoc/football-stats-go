# xG Shot Map Scraping Implementation Summary

## Overview
The football-stats-go API has been updated to scrape Expected Goals (xG) shot map data from xgstat.com. The scraper extracts detailed match statistics including shot coordinates, xG values, player names, and match outcomes.

## Changes Made

### 1. Data Models (db_models.go)
✅ **Already existed** - No changes needed
- `DBXGStatFixture`: Main fixture model with match details and shots
- `DBXGStatShot`: Individual shot data with coordinates and xG values

### 2. Scraper Service (service.go)
✅ **Updated** - Added new scraping method
- Added `ScrapeXGStatFixture()`: Main method to scrape xG data from xgstat.com
- Added `parseXGStatData()`: Helper to parse JavaScript data from the page
- Added `extractIDFromURL()`: Helper to extract match ID from URL
- Imports: Added `encoding/json`, `regexp`, `strconv`, `strings`, and domain package

**Key Features:**
- Extracts data from `window.__NUXT__` JavaScript object
- Parses match details, scores, xG values, and shot coordinates
- Handles bot detection with realistic browser headers
- 60-second timeout for page loading
- Debug mode support

### 3. API Handler (handler.go)
✅ **Updated** - Added new endpoint
- Added `ScrapeXGStats()`: HTTP handler for POST /api/scrape/xgstats
- Added domain import for DBXGStatFixture types
- Swagger documentation annotations included

### 4. Main Application (main.go)
✅ **Updated** - Registered new route
- Added route: `POST /api/scrape/xgstats`
- Updated endpoint listing in logs

### 5. Documentation
✅ **Created/Updated**

**README.md**
- Added xG Stats endpoint as primary endpoint
- Included example request/response with shot data
- Moved legacy scrape endpoint to secondary position

**XGSTATS_GUIDE.md** (New)
- Comprehensive guide for using the xG scraper
- Data model documentation
- URL pattern examples
- Testing instructions
- Troubleshooting section

### 6. Test Scripts
✅ **Created**

**test-xgstats.ps1** (PowerShell)
- Tests xG stats endpoint
- Displays formatted match summary
- Tests multiple URLs
- Color-coded output

**test-xgstats.sh** (Bash)
- Linux/Mac equivalent
- Uses jq for JSON parsing
- Formatted terminal output

### 7. Swagger Documentation
✅ **Regenerated**
- Updated with new endpoint
- Includes DBXGStatFixture and DBXGStatShot models
- Available at http://localhost:8080/swagger/

## API Endpoint

### POST /api/scrape/xgstats

**Request:**
```json
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
    "away_shots": [...]
  }
}
```

## Testing

### Start the Server
```bash
cd football-stats-go-api
make run
```

### Run Tests

**PowerShell:**
```powershell
.\test-xgstats.ps1
```

**Bash:**
```bash
chmod +x test-xgstats.sh
./test-xgstats.sh
```

**cURL:**
```bash
curl -X POST http://localhost:8080/api/scrape/xgstats \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"}'
```

## Configuration

### Environment Variables

```bash
# Run browser in visible mode (for debugging)
export SCRAPER_HEADLESS=false

# Enable debug logging
export SCRAPER_DEBUG=true
```

## Implementation Details

### Data Extraction Process

1. **Navigate to URL**: Uses chromedp to load the match page
2. **Wait for Content**: Waits for xG Shot Map section to be visible
3. **Extract Data**: Captures `window.__NUXT__` JavaScript object
4. **Parse JSON**: Parses the NUXT data structure
5. **Map to Model**: Converts to DBXGStatFixture structure
6. **Return**: Sends structured JSON response

### Anti-Bot Measures

- Realistic user agent
- Disabled automation flags
- Human-like delays
- Proper window sizing
- WebDriver property masking

### Error Handling

- 60-second timeout for page loading
- Validation of required fields
- Graceful handling of missing data
- Detailed error messages in debug mode

## Files Modified/Created

### Modified
- ✅ `internal/scraper/service.go` - Added xG scraping logic
- ✅ `internal/api/handler.go` - Added xG stats endpoint
- ✅ `cmd/api/main.go` - Registered new route
- ✅ `README.md` - Updated API documentation
- ✅ `docs/` - Regenerated Swagger docs

### Created
- ✅ `test-xgstats.ps1` - PowerShell test script
- ✅ `test-xgstats.sh` - Bash test script
- ✅ `XGSTATS_GUIDE.md` - Comprehensive usage guide

## Next Steps

Consider implementing:

1. **Database Integration**
   - Store scraped fixtures in PostgreSQL
   - Create migration for xG stats tables
   - Add CRUD endpoints for stored data

2. **Batch Scraping**
   - Scrape multiple matches in one request
   - Add queue system for async scraping
   - Schedule automatic scraping jobs

3. **Data Analysis**
   - Add endpoints to query shot patterns
   - Aggregate xG statistics by team/player
   - Generate insights from shot maps

4. **Visualization**
   - Create shot map visualization endpoints
   - Generate heatmaps of shot locations
   - Export data for frontend visualization

5. **Caching**
   - Implement Redis caching for scraped data
   - Avoid re-scraping recently fetched matches
   - Set TTL based on match status (live vs completed)

## Notes

⚠️ **Important Considerations:**

- **Rate Limiting**: Implement delays between requests to avoid overwhelming xgstat.com
- **Legal**: Ensure compliance with xgstat.com's terms of service
- **Data Accuracy**: The scraper relies on xgstat.com's data structure; changes to their site may break parsing
- **Bot Detection**: xgstat.com may implement stricter bot detection; monitor and adjust as needed

## Example URL Patterns

xgstat.com match URLs follow this pattern:
```
https://www.xgstat.com/competitions/{competition}/{season}/matches/{match-slug}/advanced-analysis/shot-maps
```

**Examples:**
- Premier League: `premier-league/2025-2026`
- La Liga: `la-liga/2025-2026`
- Bundesliga: `bundesliga/2025-2026`
- Champions League: `champions-league/2025-2026`
