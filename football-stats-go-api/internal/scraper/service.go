package scraper

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"example/hello/internal/domain"

	"github.com/chromedp/chromedp"
)

// Service provides web scraping capabilities
type Service struct {
	headless bool
	debug    bool
}

// NewService creates a new scraper service
func NewService() *Service {
	// Check environment variables for debug settings
	headless := os.Getenv("SCRAPER_HEADLESS") != "false"
	debug := os.Getenv("SCRAPER_DEBUG") == "true"

	if !headless {
		log.Println("🔍 Scraper running in VISIBLE mode - browser will be shown")
	}
	if debug {
		log.Println("🐛 Scraper debug mode enabled")
	}

	return &Service{
		headless: headless,
		debug:    debug,
	}
}

// ScrapeData represents scraped data from a website
type ScrapeData struct {
	URL       string                 `json:"url"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ScrapeXGStatFixture scrapes xG shot map data from xgstat.com
func (s *Service) ScrapeXGStatFixture(url string) (*domain.DBXGStatFixture, error) {
	if s.debug {
		log.Printf("🌐 Starting xG stat scrape for: %s", url)
	}

	// Prepare Chrome options to bypass bot detection
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
		chromedp.NoSandbox,                           // Required for Cloud Run
		chromedp.Flag("disable-dev-shm-usage", true), // Overcome limited resource problems
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("single-process", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("excludeSwitches", "enable-automation"),
		chromedp.Flag("disable-extensions", false),
		chromedp.WindowSize(1920, 1080),
	}

	// Use system Chrome if available (Cloud Run)
	if chromePath := os.Getenv("CHROME_PATH"); chromePath != "" {
		opts = append(opts, chromedp.ExecPath(chromePath))
		if s.debug {
			log.Printf("🔧 Using Chrome at: %s", chromePath)
		}
	}

	if s.headless {
		opts = append(opts, chromedp.Headless)
	} else {
		log.Printf("👀 Opening browser window for: %s", url)
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	var ctx context.Context
	var ctxCancel context.CancelFunc
	if s.debug {
		ctx, ctxCancel = chromedp.NewContext(allocCtx, chromedp.WithDebugf(log.Printf))
	} else {
		ctx, ctxCancel = chromedp.NewContext(allocCtx)
	}
	defer ctxCancel()

	ctx, timeoutCancel := context.WithTimeout(ctx, 60*time.Second)
	defer timeoutCancel()

	var pageData string
	var matchTitle string

	// Run scraping tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`, nil),
		chromedp.WaitVisible(`//h3[contains(@class, 'text-card-title') and contains(text(), 'xG Shot Map')]`, chromedp.BySearch),
		chromedp.Sleep(5*time.Second), // Wait for data to load
		chromedp.Title(&matchTitle),
		// Extract the entire page HTML content
		chromedp.OuterHTML(`html`, &pageData, chromedp.ByQuery),
	)

	if err != nil {
		log.Printf("❌ Failed to scrape %s: %v", url, err)
		return nil, fmt.Errorf("failed to scrape xG stats: %w", err)
	}

	if pageData == "" {
		return nil, fmt.Errorf("no data found on page")
	}

	// Parse the fixture data from the page data
	fixture, err := s.parseXGStatData(pageData, url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xG data: %w", err)
	}

	log.Printf("✅ Successfully scraped xG stats: %s vs %s", fixture.HomeTeam, fixture.AwayTeam)
	return fixture, nil
}

// parseXGStatData parses the raw page data into DBXGStatFixture
func (s *Service) parseXGStatData(pageData, url string) (*domain.DBXGStatFixture, error) {
	fixture := &domain.DBXGStatFixture{
		HomeShots: []domain.DBXGStatShot{},
		AwayShots: []domain.DBXGStatShot{},
	}

	// Extract team names - look for pattern like "Arsenal" and "Manchester Utd" in the header
	// Pattern: <span class="...">TeamName</span><span class="lg:hidden">ShortName</span></a><span class="font-bold"> score - score </span>
	teamPattern := regexp.MustCompile(`<span class="hidden lg:inline">([^<]+)</span><span class="lg:hidden">[^<]+</span></a><span class="font-bold">\s*(\d+)\s*-\s*(\d+)\s*</span><a[^>]*><span class="hidden lg:inline">([^<]+)</span>`)
	if teamMatches := teamPattern.FindStringSubmatch(pageData); len(teamMatches) >= 5 {
		fixture.HomeTeam = teamMatches[1]
		fixture.HomeScore, _ = strconv.Atoi(teamMatches[2])
		fixture.AwayScore, _ = strconv.Atoi(teamMatches[3])
		fixture.AwayTeam = teamMatches[4]
		if s.debug {
			log.Printf("🔍 Found teams: %s (%d) vs %s (%d)", fixture.HomeTeam, fixture.HomeScore, fixture.AwayTeam, fixture.AwayScore)
		}
	}

	// Extract xG values - look for "Expected Goals" section with values like "1.25 - 0.87"
	xgPattern := regexp.MustCompile(`<span class="tabular-nums">(\d+\.\d+)</span><span[^>]*>-</span><span class="tabular-nums">(\d+\.\d+)</span>`)
	if xgMatches := xgPattern.FindStringSubmatch(pageData); len(xgMatches) >= 3 {
		fixture.HomeXG, _ = strconv.ParseFloat(xgMatches[1], 64)
		fixture.AwayXG, _ = strconv.ParseFloat(xgMatches[2], 64)
		if s.debug {
			log.Printf("🔍 Found xG: %.2f - %.2f", fixture.HomeXG, fixture.AwayXG)
		}
	}

	// Extract gameweek - look for pattern "GW23"
	gwPattern := regexp.MustCompile(`>GW(\d+)</span>`)
	if gwMatches := gwPattern.FindStringSubmatch(pageData); len(gwMatches) >= 2 {
		fixture.Gameweek, _ = strconv.Atoi(gwMatches[1])
		if s.debug {
			log.Printf("🔍 Found gameweek: %d", fixture.Gameweek)
		}
	}

	// Extract date - look for pattern like "25 Jan 16:30"
	datePattern := regexp.MustCompile(`<span class="text-foreground text-nowrap">(\d+)\s+(\w+)\s+(\d+:\d+)</span>`)
	if dateMatches := datePattern.FindStringSubmatch(pageData); len(dateMatches) >= 4 {
		// Parse date - would need proper parsing with year, for now just log
		if s.debug {
			log.Printf("🔍 Found date: %s %s %s", dateMatches[1], dateMatches[2], dateMatches[3])
		}
	}

	// Extract ID from URL
	fixture.ID = extractIDFromURL(url)

	// Extract shot data for home team (Arsenal xG Shot Map section)
	homeMapPattern := regexp.MustCompile(`(?s)<h3[^>]*>` + regexp.QuoteMeta(fixture.HomeTeam) + ` xG Shot Map</h3>.*?</div>\s*</div>\s*</div>`)
	if homeMapMatch := homeMapPattern.FindString(pageData); homeMapMatch != "" {
		fixture.HomeShots = s.extractShotsFromSection(homeMapMatch, true)
		if s.debug {
			log.Printf("🔍 Found %d home team shots", len(fixture.HomeShots))
		}
	}

	// Extract shot data for away team (Manchester United xG Shot Map section)
	awayMapPattern := regexp.MustCompile(`(?s)<h3[^>]*>` + regexp.QuoteMeta(fixture.AwayTeam) + ` xG Shot Map</h3>.*?</div>\s*</div>\s*</div>`)
	if awayMapMatch := awayMapPattern.FindString(pageData); awayMapMatch != "" {
		fixture.AwayShots = s.extractShotsFromSection(awayMapMatch, false)
		if s.debug {
			log.Printf("🔍 Found %d away team shots", len(fixture.AwayShots))
		}
	}

	if s.debug {
		log.Printf("✅ Extracted fixture data: %s vs %s with %d total shots", fixture.HomeTeam, fixture.AwayTeam, len(fixture.HomeShots)+len(fixture.AwayShots))
	}

	return fixture, nil
}

// extractShotsFromSection extracts shot data from a team's shot map section
func (s *Service) extractShotsFromSection(sectionHTML string, isHomeTeam bool) []domain.DBXGStatShot {
	shots := []domain.DBXGStatShot{}

	// Build a map of player data from the table
	playerData := make(map[string]*domain.DBXGStatShot)

	// Extract player rows from the table
	// Pattern: player number, name, xG value, goals, shots on target
	playerRowPattern := regexp.MustCompile(`(?s)<span class="flex size-4[^>]*>(\d+)</span><span[^>]*>([^<]+)</span>.*?<div class="rounded[^>]*>(\d+\.?\d*)</div>.*?<div class="rounded[^>]*>(\d+)</div>`)
	playerMatches := playerRowPattern.FindAllStringSubmatch(sectionHTML, -1)

	for _, match := range playerMatches {
		if len(match) >= 5 {
			playerName := match[2]
			xgValue, _ := strconv.ParseFloat(match[3], 64)
			goals, _ := strconv.Atoi(match[4])

			shot := &domain.DBXGStatShot{
				PlayerName: playerName,
				XG:         xgValue,
				IsGoal:     goals > 0,
				ShotType:   "", // Will be determined from circle class
			}
			playerData[playerName] = shot
		}
	}

	// Extract SVG circle data for shot positions
	// Different types: off-target (fill="var(--foreground)" fill-opacity="0.3"),
	// blocked (fill="var(--chart-red)"), on-target (fill-opacity="0.9"), goals (star svg)

	// Off-target shots
	offTargetPattern := regexp.MustCompile(`<circle r="[\d.]+" cx="([\d.]+)" cy="([\d.]+)"[^>]*fill="var\(--foreground\)" fill-opacity="0\.3"`)
	for _, match := range offTargetPattern.FindAllStringSubmatch(sectionHTML, -1) {
		if len(match) >= 3 {
			cx, _ := strconv.ParseFloat(match[1], 64)
			cy, _ := strconv.ParseFloat(match[2], 64)
			shots = append(shots, domain.DBXGStatShot{
				X:        cx,
				Y:        cy,
				ShotType: "off_target",
			})
		}
	}

	// Blocked shots
	blockedPattern := regexp.MustCompile(`<circle r="[\d.]+" cx="([\d.]+)" cy="([\d.]+)"[^>]*fill="var\(--chart-red\)"`)
	for _, match := range blockedPattern.FindAllStringSubmatch(sectionHTML, -1) {
		if len(match) >= 3 {
			cx, _ := strconv.ParseFloat(match[1], 64)
			cy, _ := strconv.ParseFloat(match[2], 64)
			shots = append(shots, domain.DBXGStatShot{
				X:        cx,
				Y:        cy,
				ShotType: "blocked",
			})
		}
	}

	// On-target shots (non-goals)
	onTargetPattern := regexp.MustCompile(`<circle r="[\d.]+" cx="([\d.]+)" cy="([\d.]+)"[^>]*fill="var\(--foreground\)" fill-opacity="0\.9"`)
	for _, match := range onTargetPattern.FindAllStringSubmatch(sectionHTML, -1) {
		if len(match) >= 3 {
			cx, _ := strconv.ParseFloat(match[1], 64)
			cy, _ := strconv.ParseFloat(match[2], 64)
			shots = append(shots, domain.DBXGStatShot{
				X:        cx,
				Y:        cy,
				ShotType: "on_target",
			})
		}
	}

	// Goals (star SVG elements)
	goalPattern := regexp.MustCompile(`<svg[^>]*x="([\d.]+)" y="([\d.]+)"[^>]*fill="var\(--brand-yellow\)"`)
	for _, match := range goalPattern.FindAllStringSubmatch(sectionHTML, -1) {
		if len(match) >= 3 {
			x, _ := strconv.ParseFloat(match[1], 64)
			y, _ := strconv.ParseFloat(match[2], 64)
			shots = append(shots, domain.DBXGStatShot{
				X:        x,
				Y:        y,
				IsGoal:   true,
				ShotType: "goal",
			})
		}
	}

	if s.debug {
		log.Printf("🔍 Extracted %d shots from section", len(shots))
	}

	return shots
}

// extractIDFromURL extracts a numeric ID from the URL
func extractIDFromURL(url string) int {
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllString(url, -1)
	if len(matches) > 0 {
		// Use the last number found in URL as ID
		for i := len(matches) - 1; i >= 0; i-- {
			if id, err := strconv.Atoi(matches[i]); err == nil {
				return id
			}
		}
	}
	return 0
}
