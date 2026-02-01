# VS Code Debugging Guide for Go

## Quick Start

### 1️⃣ Set a Breakpoint
- Open any `.go` file (e.g., `football-stats-go-api/cmd/api/main.go` or `internal/scraper/service.go`)
- Click in the **left margin** (gutter) next to the line number where you want to pause
- A **red dot** will appear - that's your breakpoint!

### 2️⃣ Start Debugging
- Press **F5** OR
- Click **Run → Start Debugging** OR  
- Click the **Run and Debug** icon in the left sidebar (play button with bug), then click the green play button

### 3️⃣ Choose Configuration
Select **"Debug API Server"** from the dropdown (runs with visible browser + debug logs)

### 4️⃣ Trigger Your Breakpoint
Once the debugger is running, make an API request:
```powershell
curl -X POST http://localhost:8080/api/scrape -H "Content-Type: application/json" -d '{"url":"https://example.com"}'
```

### 5️⃣ Debug Controls
When stopped at a breakpoint:
- **F10** - Step Over (next line)
- **F11** - Step Into (enter function)
- **Shift+F11** - Step Out (exit function)
- **F5** - Continue (run to next breakpoint)
- **Shift+F5** - Stop Debugging

## Good Places to Set Breakpoints

### In Handler ([internal/api/handler.go](../football-stats-go-api/internal/api/handler.go))
```go
func (h *Handler) Scrape(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    var req ScrapeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        // ⭕ SET BREAKPOINT HERE to debug request parsing
        writeError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // ⭕ SET BREAKPOINT HERE to see the URL being scraped
    data, err := h.scraperService.ScrapeWebsite(req.URL)
    if err != nil {
        // ⭕ SET BREAKPOINT HERE to catch scraper errors
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    writeSuccess(w, data)
}
```

### In Scraper Service ([internal/scraper/service.go](../football-stats-go-api/internal/scraper/service.go))
```go
func (s *Service) ScrapeWebsite(url string) (*ScrapeData, error) {
    // ⭕ SET BREAKPOINT HERE to start debugging the scrape
    if s.debug {
        log.Printf("🌐 Starting scrape for: %s", url)
    }

    // ... browser setup ...

    // ⭕ SET BREAKPOINT HERE before navigation
    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.WaitVisible(`body`, chromedp.ByQuery),
        // ⭕ SET BREAKPOINT in this Run block to debug scraping steps
        chromedp.Title(&title),
        chromedp.Text(`body`, &content, chromedp.ByQuery),
    )

    if err != nil {
        // ⭕ SET BREAKPOINT HERE to catch scraping errors
        return nil, fmt.Errorf("failed to scrape website: %w", err)
    }

    // ⭕ SET BREAKPOINT HERE to inspect scraped data
    return &ScrapeData{
        URL:       url,
        Title:     title,
        Content:   content,
        // ...
    }, nil
}
```

## Debug Panel Features

When debugging, the **left sidebar** shows:
- **Variables** - Inspect all local variables, function parameters
- **Watch** - Add expressions to monitor (e.g., `len(content)`, `req.URL`)
- **Call Stack** - See the function call hierarchy
- **Breakpoints** - List of all your breakpoints

## Tips

✅ **Inspect Variables**: Hover over any variable to see its value  
✅ **Debug Console**: Type Go expressions to evaluate while paused  
✅ **Conditional Breakpoints**: Right-click red dot → Edit Breakpoint → Add condition (e.g., `url == "https://example.com"`)  
✅ **Logpoints**: Right-click gutter → Add Logpoint (logs without stopping)  
✅ **Multiple Breakpoints**: Set as many as you need!  

## Available Configurations

1. **Debug API Server** - Visible browser + debug logs (recommended for development)
2. **Debug API Server (Headless)** - Hidden browser + debug logs  
3. **Debug API Server (Production Mode)** - Production environment simulation

Switch between them in the **Run and Debug** dropdown.

## Install Go Extension (if needed)

If debugging doesn't work, install the Go extension:
1. Press **Ctrl+Shift+X** (Extensions)
2. Search for "Go"
3. Install the official **Go** extension by Go Team at Google
4. Reload VS Code

The extension will automatically install `delve` (Go debugger) when you first start debugging.
