# Debugging the Scraper

## Running with Visible Browser

To see the browser window during scraping, set the environment variable:

### PowerShell
```powershell
$env:SCRAPER_HEADLESS="false"
go run cmd/api/main.go
```

### Bash/Linux
```bash
SCRAPER_HEADLESS=false go run cmd/api/main.go
```

## Enable Debug Logging

For detailed step-by-step logging:

### PowerShell
```powershell
$env:SCRAPER_DEBUG="true"
go run cmd/api/main.go
```

### Both Together
```powershell
$env:SCRAPER_HEADLESS="false"
$env:SCRAPER_DEBUG="true"
go run cmd/api/main.go
```

## Test the Scraper

Once the server is running, test it:

```powershell
# Simple test
Invoke-RestMethod -Method POST -Uri "http://localhost:8080/api/scrape" `
  -ContentType "application/json" `
  -Body '{"url":"https://example.com"}' | ConvertTo-Json -Depth 5

# Or with curl
curl -X POST http://localhost:8080/api/scrape `
  -H "Content-Type: application/json" `
  -d '{"url":"https://example.com"}'
```

## Debug Output

With debugging enabled, you'll see:
- 🌐 Starting scrape for URL
- 👀 Opening browser window (if not headless)
- 📍 Navigating to URL
- ⏳ Waiting for page to load
- 📄 Extracting page data
- ✅ Successfully scraped with stats
- ❌ Error messages if something fails

## Common Issues

### Chrome/Chromium Not Found
Install Chrome or Chromium:
```powershell
# Using Chocolatey
choco install googlechrome
```

### Timeout Errors
Increase timeout in `internal/scraper/service.go`:
```go
ctx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second) // Increase from 30s
```

### Page Not Loading
Try adding a sleep to wait for dynamic content:
```go
chromedp.Sleep(2 * time.Second)
```

## Environment Variables Summary

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `SCRAPER_HEADLESS` | `true`/`false` | `true` | Show browser window when false |
| `SCRAPER_DEBUG` | `true`/`false` | `false` | Enable verbose logging |

## VS Code Launch Configuration

Add to `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Scraper (Visible)",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/football-stats-go-api/cmd/api",
      "env": {
        "SCRAPER_HEADLESS": "false",
        "SCRAPER_DEBUG": "true"
      }
    }
  ]
}
```
