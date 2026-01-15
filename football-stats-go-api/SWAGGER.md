# Swagger Documentation Setup

## Quick Start

### 1. Install Swag CLI Tool

```bash
# Using make
make install-swag

# Or manually
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Generate Swagger Documentation

```bash
# Using make (recommended)
make swagger

# Or manually
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### 3. Run the API

```bash
# Using make (auto-generates docs)
make run

# Or manually
go run cmd/api/main.go
```

### 4. Access Swagger UI

Open your browser to:
```
http://localhost:8080/swagger/
```

## Swagger Annotations

The API uses **swaggo/swag** annotations in code comments to generate OpenAPI/Swagger documentation.

### Example Handler Annotation

```go
// AnalyzeTeam analyzes team statistics
// @Summary Analyze team statistics
// @Description Calculate comprehensive statistics for a team based on match data
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body AnalyzeTeamRequest true "Team analysis request"
// @Success 200 {object} Response{data=domain.TeamStats} "Team statistics"
// @Failure 400 {object} Response "Invalid request"
// @Router /analyze-team [post]
func (h *Handler) AnalyzeTeam(w http.ResponseWriter, r *http.Request) {
    // handler code
}
```

### General API Info (in handler.go)

```go
// @title Football Stats Analysis API
// @version 1.0
// @description A lightweight API for analyzing football data
// @host localhost:8080
// @BasePath /api
```

## Available Endpoints in Swagger

- **POST /api/analyze-team** - Analyze team statistics
- **POST /api/predict-match** - Predict match outcome
- **POST /api/head-to-head** - Head-to-head analysis
- **GET /health** - Health check

## Updating Documentation

After changing handler annotations:

```bash
# Regenerate docs
make swagger

# Or
swag init -g cmd/api/main.go -o docs
```

The generated files are in `docs/` directory:
- `docs/docs.go` - Go code
- `docs/swagger.json` - OpenAPI JSON spec
- `docs/swagger.yaml` - OpenAPI YAML spec

## CI/CD Integration

The Dockerfile already includes Swagger generation:

```dockerfile
# In your build process
RUN swag init -g cmd/api/main.go -o docs
```

For GitHub Actions, add before build:

```yaml
- name: Generate Swagger docs
  run: |
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g cmd/api/main.go -o docs
```

## Alternative Tools

### 1. **swaggo/swag** (Current - Recommended ✅)
- **Pros**: Code-first, lightweight, popular, easy annotations
- **Cons**: Need to regenerate docs after changes
- **Best for**: Your use case - simple API with Go-native approach

### 2. **go-swagger**
- **Pros**: More powerful, code generation from specs
- **Cons**: Heavier, more complex
- **Best for**: Large enterprise APIs with spec-first approach

### 3. **oapi-codegen**
- **Pros**: Generates server/client from OpenAPI specs
- **Cons**: Spec-first (opposite approach)
- **Best for**: When you have OpenAPI spec first

### 4. **Manual OpenAPI YAML**
- **Pros**: Full control, no dependencies
- **Cons**: Manual maintenance, prone to errors
- **Best for**: Simple APIs, documentation-focused projects

## Swagger vs Alternatives

| Tool | Approach | Learning Curve | Maintenance |
|------|----------|----------------|-------------|
| **swaggo/swag** ✅ | Code annotations | Low | Auto-generated |
| **go-swagger** | Spec-first | High | Complex |
| **oapi-codegen** | Spec-first | Medium | Spec-driven |
| **Postman** | Manual collection | Low | Manual |
| **Insomnia** | Manual collection | Low | Manual |

**Recommendation**: Stick with **swaggo/swag** for your lightweight API. It's perfect for Go projects and widely adopted.

## Tips

1. **Keep annotations updated**: Run `make swagger` after changing handlers
2. **Use tags**: Organize endpoints with `@Tags` for better UI
3. **Document models**: Add JSON tags and comments to structs
4. **Version your API**: Use `@version` annotation
5. **Test examples**: Add `@Param` examples for better docs

## Troubleshooting

### Swagger UI shows "Failed to load API definition"
- Run `make swagger` to regenerate docs
- Check for syntax errors in annotations

### Import cycle error
- Use `--parseInternal` flag with swag init

### Docs not updating
- Delete `docs/` folder and regenerate
- Restart the server
