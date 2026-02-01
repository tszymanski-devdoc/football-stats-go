# Quick Reference - xG Shot Map Scraping

## 🚀 Quick Start

### 1. Start the Server
```bash
cd football-stats-go-api
make run
```
Server will be available at `http://localhost:8080`

### 2. Test the Endpoint

**Using PowerShell:**
```powershell
.\test-xgstats.ps1
```

**Using Bash:**
```bash
chmod +x test-xgstats.sh
./test-xgstats.sh
```

**Using cURL:**
```bash
curl -X POST http://localhost:8080/api/scrape/xgstats \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"
  }'
```

## 📡 API Endpoint

```
POST /api/scrape/xgstats
```

### Request Body
```json
{
  "url": "<xgstat.com match URL>"
}
```

### Response
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
    "home_shots": [...],
    "away_shots": [...]
  }
}
```

## 🔧 Environment Variables

```bash
# Show browser (for debugging)
export SCRAPER_HEADLESS=false

# Enable debug logging
export SCRAPER_DEBUG=true
```

## 📊 Data Structure

### Shot Object
```json
{
  "x": 88.5,           // X coordinate (0-100)
  "y": 45.2,           // Y coordinate (0-100)
  "xg": 0.45,          // Expected goals (0-1)
  "is_goal": true,     // Whether it was a goal
  "shot_type": "Right Foot",
  "player_name": "Bukayo Saka",
  "minute": 23
}
```

## 🌐 Finding Match URLs

Pattern:
```
https://www.xgstat.com/competitions/{league}/{season}/matches/{match-slug}/advanced-analysis/shot-maps
```

Examples:
- `premier-league/2025-2026`
- `la-liga/2025-2026`
- `champions-league/2025-2026`

## 📚 Documentation Files

- `XGSTATS_GUIDE.md` - Comprehensive usage guide
- `XGSTATS_IMPLEMENTATION.md` - Implementation details
- `README.md` - Updated with xG endpoint
- Swagger UI: http://localhost:8080/swagger/

## 🧪 Testing

| Script | Platform | Features |
|--------|----------|----------|
| `test-xgstats.ps1` | Windows | Color output, multiple tests |
| `test-xgstats.sh` | Linux/Mac | jq formatting, color output |

## ⚠️ Important Notes

1. **Rate Limiting**: Add delays between requests
2. **Legal**: Check xgstat.com terms of service
3. **Bot Detection**: May need updates if site changes
4. **Data Structure**: Parsing depends on current site structure

## 📂 Files Modified

- ✅ `internal/scraper/service.go` - Main scraping logic
- ✅ `internal/api/handler.go` - API endpoint
- ✅ `cmd/api/main.go` - Route registration
- ✅ `README.md` - Documentation
- ✅ `docs/*` - Swagger docs

## 📂 Files Created

- ✅ `test-xgstats.ps1` - PowerShell test
- ✅ `test-xgstats.sh` - Bash test
- ✅ `XGSTATS_GUIDE.md` - Usage guide
- ✅ `XGSTATS_IMPLEMENTATION.md` - Implementation details
- ✅ `QUICK_REFERENCE.md` - This file

## 🔍 Troubleshooting

### Build Fails
```bash
cd football-stats-go-api
go mod tidy
make build
```

### Scraper Times Out
- Increase timeout in `service.go` (line 90)
- Check URL is correct
- Run in visible mode: `SCRAPER_HEADLESS=false`

### No Data Returned
- Enable debug: `SCRAPER_DEBUG=true`
- Check if match exists on xgstat.com
- Verify data structure hasn't changed

### Bot Detection
- Try non-headless mode first
- Add longer delays
- Check xgstat.com accessibility

## 🎯 Next Steps

1. ✅ Build and test the scraper
2. 🔲 Add database storage for fixtures
3. 🔲 Create batch scraping functionality
4. 🔲 Add caching layer
5. 🔲 Implement rate limiting
6. 🔲 Create data analysis endpoints
