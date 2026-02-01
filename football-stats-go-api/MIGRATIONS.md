# Database Migrations Guide

## Overview

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations. Migrations are stored in the `migrations/` folder as SQL files.

## Database Schema

The database consists of two main tables for storing xG (expected goals) statistics:

### Tables

1. **xgstat_fixtures** - Stores match fixture data
   - Basic match information (teams, scores, date, gameweek)
   - Expected goals (xG) for home and away teams
   - Indexed on gameweek and fixture_date for faster queries

2. **xgstat_shots** - Stores individual shot data
   - Shot coordinates (x, y)
   - xG value for each shot
   - Player information and shot type
   - Links to fixtures via foreign key
   - Indexed on fixture_id, player_name, and is_goal

## Prerequisites

1. **PostgreSQL Database** - Running PostgreSQL instance
2. **Go 1.24+** - For running the migration tool
3. **Database URL** - Connection string to your PostgreSQL database

## Local Development

### 1. Set Database Connection

Set the `DATABASE_URL` environment variable:

**PowerShell:**
```powershell
$env:DATABASE_URL = "postgresql://username:password@localhost:5432/dbname?sslmode=disable"
```

**Bash:**
```bash
export DATABASE_URL="postgresql://username:password@localhost:5432/dbname?sslmode=disable"
```

### 2. Run Migrations

**Apply all migrations (up):**
```bash
make migrate-up
```

Or directly:
```bash
go run cmd/migrate/main.go -direction up
```

**Rollback migrations (down):**
```bash
make migrate-down
```

Or directly:
```bash
go run cmd/migrate/main.go -direction down
```

**With custom database URL:**
```bash
go run cmd/migrate/main.go -direction up -db "postgresql://user:pass@host:5432/db?sslmode=disable"
```

## Production/Cloud Deployment

Migrations are automatically run as part of the CD pipeline before deploying the application to Cloud Run.

### CD Pipeline Migration Job

The `.github/workflows/api-cd.yml` includes a `migrate` job that:

1. Sets up Go environment
2. Authenticates with Google Cloud
3. Starts Cloud SQL Proxy to connect to the database
4. Runs migrations using the migrate command
5. Deploys the application only if migrations succeed

### Required Secrets

Configure these secrets in your GitHub repository:

- `GCP_SA_KEY` - Google Cloud service account key JSON
- `GCP_PROJECT_ID` - Your GCP project ID
- `GCP_REGION` - GCP region where your database is located (e.g., `us-central1`)
- `DB_USER` - Database username for football-stats-go-db
- `DB_PASSWORD` - Database password for football-stats-go-db

**Note:** The pipeline uses the existing `football-stats-go-db` PostgreSQL 18 instance in GCP and connects to the default `postgres` database.

## Migration Files Structure

```
migrations/
├── 000001_create_xgstat_tables.up.sql    # Creates tables and indexes
└── 000001_create_xgstat_tables.down.sql  # Drops tables
```

## Creating New Migrations

To create a new migration:

1. Create two files in the `migrations/` folder with the next version number:
   - `XXXXXX_description.up.sql` - Forward migration
   - `XXXXXX_description.down.sql` - Rollback migration

2. Version numbers should be sequential (e.g., 000001, 000002, etc.)

Example:
```sql
-- migrations/000002_add_column_to_fixtures.up.sql
ALTER TABLE xgstat_fixtures ADD COLUMN competition VARCHAR(100);

-- migrations/000002_add_column_to_fixtures.down.sql
ALTER TABLE xgstat_fixtures DROP COLUMN competition;
```

## Troubleshooting

### Migration fails with "dirty database"

If a migration fails partway through, the database may be marked as "dirty". To resolve:

1. Check the current migration version:
   ```sql
   SELECT * FROM schema_migrations;
   ```

2. Manually fix the issue and mark as clean, or force to a specific version using the golang-migrate CLI

### Connection issues

Ensure:
- PostgreSQL is running
- Database credentials are correct
- Network connectivity is available
- SSL mode is appropriate for your environment

### Permission errors

Make sure the database user has sufficient permissions:
```sql
GRANT ALL PRIVILEGES ON DATABASE your_database TO your_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_user;
```

## Best Practices

1. **Always create both up and down migrations** - Allows rollback if needed
2. **Test migrations locally first** - Before deploying to production
3. **Keep migrations small and focused** - One logical change per migration
4. **Never edit existing migrations** - Create new ones instead
5. **Use transactions where possible** - But note some DDL statements can't be rolled back in PostgreSQL
6. **Back up production data** - Before running migrations in production

## Manual Database Setup (Alternative)

If you prefer to set up the database manually without migrations:

```sql
-- See migrations/000001_create_xgstat_tables.up.sql for the full schema
```

However, using migrations is recommended for consistency and version control.
