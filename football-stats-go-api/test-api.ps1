# Test the Football Analysis API (PowerShell)

Write-Host "Testing Football Analysis API" -ForegroundColor Cyan

# Health check
Write-Host "`n1. Health Check" -ForegroundColor Yellow
curl http://localhost:8080/health

# Analyze team statistics
Write-Host "`n2. Analyze Team Statistics" -ForegroundColor Yellow
$analyzeBody = @{
    team_name = "Manchester United"
    matches = @(
        @{
            id = "1"
            home_team = "Manchester United"
            away_team = "Liverpool"
            home_score = 2
            away_score = 1
            match_date = "2024-01-15T15:00:00Z"
            league = "Premier League"
            season = "2024"
        },
        @{
            id = "2"
            home_team = "Chelsea"
            away_team = "Manchester United"
            home_score = 1
            away_score = 1
            match_date = "2024-01-22T15:00:00Z"
            league = "Premier League"
            season = "2024"
        },
        @{
            id = "3"
            home_team = "Manchester United"
            away_team = "Arsenal"
            home_score = 3
            away_score = 0
            match_date = "2024-01-29T15:00:00Z"
            league = "Premier League"
            season = "2024"
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/api/analyze-team" `
    -Method Post `
    -ContentType "application/json" `
    -Body $analyzeBody | ConvertTo-Json

# Predict match
Write-Host "`n3. Predict Match Outcome" -ForegroundColor Yellow
$predictBody = @{
    home_team = "Manchester United"
    away_team = "Liverpool"
    home_matches = @(
        @{
            id = "1"
            home_team = "Manchester United"
            away_team = "Chelsea"
            home_score = 2
            away_score = 1
            match_date = "2024-01-15T15:00:00Z"
            league = "Premier League"
            season = "2024"
        },
        @{
            id = "2"
            home_team = "Manchester United"
            away_team = "Arsenal"
            home_score = 3
            away_score = 0
            match_date = "2024-01-22T15:00:00Z"
            league = "Premier League"
            season = "2024"
        }
    )
    away_matches = @(
        @{
            id = "3"
            home_team = "Liverpool"
            away_team = "Chelsea"
            home_score = 4
            away_score = 1
            match_date = "2024-01-15T15:00:00Z"
            league = "Premier League"
            season = "2024"
        },
        @{
            id = "4"
            home_team = "Liverpool"
            away_team = "Arsenal"
            home_score = 2
            away_score = 2
            match_date = "2024-01-22T15:00:00Z"
            league = "Premier League"
            season = "2024"
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/api/predict-match" `
    -Method Post `
    -ContentType "application/json" `
    -Body $predictBody | ConvertTo-Json

# Head-to-head
Write-Host "`n4. Head-to-Head Analysis" -ForegroundColor Yellow
$h2hBody = @{
    team1 = "Manchester United"
    team2 = "Liverpool"
    matches = @(
        @{
            id = "1"
            home_team = "Manchester United"
            away_team = "Liverpool"
            home_score = 2
            away_score = 1
            match_date = "2023-09-15T15:00:00Z"
            league = "Premier League"
            season = "2023"
        },
        @{
            id = "2"
            home_team = "Liverpool"
            away_team = "Manchester United"
            home_score = 3
            away_score = 1
            match_date = "2023-12-10T15:00:00Z"
            league = "Premier League"
            season = "2023"
        },
        @{
            id = "3"
            home_team = "Manchester United"
            away_team = "Liverpool"
            home_score = 1
            away_score = 1
            match_date = "2024-03-20T15:00:00Z"
            league = "Premier League"
            season = "2024"
        }
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Uri "http://localhost:8080/api/head-to-head" `
    -Method Post `
    -ContentType "application/json" `
    -Body $h2hBody | ConvertTo-Json

Write-Host "`nAll tests completed!" -ForegroundColor Green
