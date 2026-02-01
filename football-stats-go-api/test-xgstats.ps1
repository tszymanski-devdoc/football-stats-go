# Test script for xG Stats scraping endpoint
Write-Host "Testing xG Stats Scraper Endpoint..." -ForegroundColor Cyan

$url = "http://localhost:8080/api/scrape/xgstats"

# Example URL from xgstat.com
$body = @{
    url = "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"
} | ConvertTo-Json

Write-Host "`nSending POST request to: $url" -ForegroundColor Yellow
Write-Host "Request Body:" -ForegroundColor Yellow
Write-Host $body -ForegroundColor Gray

try {
    $response = Invoke-RestMethod -Uri $url -Method Post -Body $body -ContentType "application/json"
    
    Write-Host "`n✅ Success!" -ForegroundColor Green
    Write-Host "`nResponse:" -ForegroundColor Cyan
    $response | ConvertTo-Json -Depth 10
    
    if ($response.success) {
        $data = $response.data
        Write-Host "`n📊 Match Summary:" -ForegroundColor Magenta
        Write-Host "  $($data.home_team) $($data.home_score) - $($data.away_score) $($data.away_team)" -ForegroundColor White
        Write-Host "  xG: $($data.home_xg) - $($data.away_xg)" -ForegroundColor Yellow
        Write-Host "  Date: $($data.date)" -ForegroundColor Gray
        Write-Host "  Gameweek: $($data.gameweek)" -ForegroundColor Gray
        Write-Host "`n  Home Shots: $($data.home_shots.Count)" -ForegroundColor Cyan
        Write-Host "  Away Shots: $($data.away_shots.Count)" -ForegroundColor Cyan
    }
} catch {
    Write-Host "`n❌ Error!" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        Write-Host $_.ErrorDetails.Message -ForegroundColor Red
    }
}

Write-Host "`n---" -ForegroundColor Gray

# Test with another match (update this URL as needed)
Write-Host "`nTesting with another match URL..." -ForegroundColor Cyan

$body2 = @{
    url = "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/liverpool-chelsea-2026-01-25/advanced-analysis/shot-maps"
} | ConvertTo-Json

try {
    $response2 = Invoke-RestMethod -Uri $url -Method Post -Body $body2 -ContentType "application/json"
    
    Write-Host "✅ Success!" -ForegroundColor Green
    
    if ($response2.success) {
        $data2 = $response2.data
        Write-Host "`n📊 Match Summary:" -ForegroundColor Magenta
        Write-Host "  $($data2.home_team) $($data2.home_score) - $($data2.away_score) $($data2.away_team)" -ForegroundColor White
        Write-Host "  xG: $($data2.home_xg) - $($data2.away_xg)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "⚠️  Second test failed (URL may not exist): $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host "`nTest complete!" -ForegroundColor Green
