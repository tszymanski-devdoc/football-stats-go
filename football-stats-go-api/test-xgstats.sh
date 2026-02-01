#!/bin/bash

# Test script for xG Stats scraping endpoint
echo -e "\033[0;36mTesting xG Stats Scraper Endpoint...\033[0m"

URL="http://localhost:8080/api/scrape/xgstats"

# Example URL from xgstat.com
BODY='{
  "url": "https://www.xgstat.com/competitions/premier-league/2025-2026/matches/arsenal-manchester-united-2026-01-24/advanced-analysis/shot-maps"
}'

echo -e "\n\033[0;33mSending POST request to: $URL\033[0m"
echo -e "\033[0;33mRequest Body:\033[0m"
echo -e "\033[0;37m$BODY\033[0m"

RESPONSE=$(curl -s -X POST "$URL" \
  -H "Content-Type: application/json" \
  -d "$BODY")

if [ $? -eq 0 ]; then
    echo -e "\n\033[0;32m✅ Success!\033[0m"
    echo -e "\n\033[0;36mResponse:\033[0m"
    echo "$RESPONSE" | jq '.'
    
    # Extract and display match summary
    HOME_TEAM=$(echo "$RESPONSE" | jq -r '.data.home_team')
    AWAY_TEAM=$(echo "$RESPONSE" | jq -r '.data.away_team')
    HOME_SCORE=$(echo "$RESPONSE" | jq -r '.data.home_score')
    AWAY_SCORE=$(echo "$RESPONSE" | jq -r '.data.away_score')
    HOME_XG=$(echo "$RESPONSE" | jq -r '.data.home_xg')
    AWAY_XG=$(echo "$RESPONSE" | jq -r '.data.away_xg')
    GAMEWEEK=$(echo "$RESPONSE" | jq -r '.data.gameweek')
    HOME_SHOTS=$(echo "$RESPONSE" | jq -r '.data.home_shots | length')
    AWAY_SHOTS=$(echo "$RESPONSE" | jq -r '.data.away_shots | length')
    
    if [ "$HOME_TEAM" != "null" ]; then
        echo -e "\n\033[0;35m📊 Match Summary:\033[0m"
        echo -e "  \033[1;37m$HOME_TEAM $HOME_SCORE - $AWAY_SCORE $AWAY_TEAM\033[0m"
        echo -e "  \033[0;33mxG: $HOME_XG - $AWAY_XG\033[0m"
        echo -e "  \033[0;37mGameweek: $GAMEWEEK\033[0m"
        echo -e "\n  \033[0;36mHome Shots: $HOME_SHOTS\033[0m"
        echo -e "  \033[0;36mAway Shots: $AWAY_SHOTS\033[0m"
    fi
else
    echo -e "\n\033[0;31m❌ Error!\033[0m"
    echo "$RESPONSE"
fi

echo -e "\n\033[0;90m---\033[0m"
echo -e "\n\033[0;32mTest complete!\033[0m"
