# Database Service Guide

## Overview

The database service handles persistence of xG statistics data scraped from xgstat.com. It uses PostgreSQL and provides automatic saving when scraping fixtures.

## Architecture

### Database Service (`internal/database/service.go`)

The service provides:
- **Connection management** - Handles database connections with pooling
- **Transaction support** - Ensures data consistency when saving fixtures and shots
- **CRUD operations** - Save and retrieve xG statistics

### Key Methods

#### `SaveXGStatFixture(fixture *domain.DBXGStatFixture) error`
Saves a complete fixture with all shots to the database. Features:
- Uses transactions for atomicity
- Updates existing fixtures (upsert based on fixture_id + gameweek)
- Deletes old shots and inserts new ones to avoid duplicates
- Saves home and away shots with team type markers

#### `GetFixtureByID(fixtureID int) (*domain.DBXGStatFixture, error)`
Retrieves a saved fixture with all associated shots from the database.

## API Endpoints

### POST /api/scrape/xgstats
Scrapes xG data from xgstat.com **and automatically saves it to the database**.

**Request:**
```json
{
  "url": "https://xgstat.com/fixture/..."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "gameweek": 1,
    "id": 12345,
    "date": "2026-02-01T15:00:00Z",
    "home_team": "Team A",
    "away_team": "Team B",
    "home_score": 2,
    "away_score": 1,
    "home_xg": 1.85,
    "away_xg": 0.92,
    "home_shots": [...],
    "away_shots": [...]
  }
}
```

### GET /api/xgstats?id=12345
Retrieves a previously saved fixture from the database.

**Query Parameters:**
- `id` (required) - The fixture ID

**Response:** Same format as scrape endpoint

## Configuration

The service supports two configuration methods:

### 1. DATABASE_URL (Recommended for Cloud/Production)
```bash
export DATABASE_URL="postgresql://user:password@host:5432/dbname?sslmode=disable"
```

### 2. Individual Environment Variables (Local Development)
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=postgres
export DB_SSL_MODE=disable
```

### Connection Pool Settings (Optional)
```bash
export DB_MAX_OPEN_CONNS=25
export DB_MAX_IDLE_CONNS=25
export DB_CONN_MAX_LIFETIME=30m
```

## Database Schema

### Tables

#### xgstat_fixtures
Stores match fixture information:
- `id` - Auto-incrementing primary key
- `gameweek` - Premier League gameweek number
- `fixture_id` - External fixture ID from xgstat.com
- `fixture_date` - Match date and time
- `home_team`, `away_team` - Team names
- `home_score`, `away_score` - Final scores
- `home_xg`, `away_xg` - Expected goals values
- `created_at`, `updated_at` - Timestamps

**Unique constraint:** `(fixture_id, gameweek)`

#### xgstat_shots
Stores individual shot data:
- `id` - Auto-incrementing primary key
- `fixture_id` - Foreign key to xgstat_fixtures
- `x`, `y` - Shot coordinates (0-100)
- `xg` - Expected goal value (0-1)
- `is_goal` - Whether the shot resulted in a goal
- `shot_type` - Type of shot (e.g., "Right foot")
- `player_name` - Name of the player
- `minute` - Match minute (1-120)
- `team_type` - Either "home" or "away"
- `created_at` - Timestamp

## Usage Examples

### Scraping and Saving
```bash
curl -X POST http://localhost:8080/api/scrape/xgstats \
  -H "Content-Type: application/json" \
  -d '{"url": "https://xgstat.com/fixture/12345"}'
```

This will:
1. Scrape the fixture data
2. Save to database (or update if exists)
3. Return the complete fixture data

### Retrieving Saved Data
```bash
curl http://localhost:8080/api/xgstats?id=12345
```

## Error Handling

The service handles various error scenarios:
- **Connection failures** - Returns error if database is unreachable
- **Duplicate prevention** - Upserts fixtures, replaces shots
- **Transaction rollback** - Ensures data consistency on failures
- **Not found** - Returns 404 if fixture doesn't exist

## Local Testing

1. Ensure PostgreSQL is running
2. Run migrations: `go run cmd/migrate/main.go -direction up`
3. Set DATABASE_URL or individual DB env vars
4. Start the API: `go run cmd/api/main.go`
5. Scrape a fixture - it will automatically save to DB
6. Retrieve it using the GET endpoint

## Production Deployment

In production (Cloud Run):
- Database connection is established via Cloud SQL Unix socket
- DATABASE_URL is set automatically by the deployment pipeline
- Connection pooling is configured for optimal performance
- Migrations run automatically before deployment

## Notes

- The scrape endpoint **always saves** to the database after successful scraping
- Existing fixtures are updated (not duplicated) based on fixture_id + gameweek
- All shots are replaced on update to ensure data accuracy
- Transactions ensure either all data is saved or none (no partial saves)
