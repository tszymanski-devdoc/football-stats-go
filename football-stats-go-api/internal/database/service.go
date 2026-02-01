package database

import (
	"database/sql"
	"fmt"

	"example/hello/internal/config"
	"example/hello/internal/domain"

	_ "github.com/lib/pq"
)

// Service handles database operations
type Service struct {
	db *sql.DB
}

// NewService creates a new database service
func NewService(cfg *config.Config) (*Service, error) {
	// Build connection string from config or use DATABASE_URL
	var connStr string
	if dbURL := cfg.Database.URL; dbURL != "" {
		connStr = dbURL
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Service{db: db}, nil
}

// Close closes the database connection
func (s *Service) Close() error {
	return s.db.Close()
}

// SaveXGStatFixture saves a fixture and its shots to the database
func (s *Service) SaveXGStatFixture(fixture *domain.DBXGStatFixture) error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert or update fixture
	var fixtureID int
	err = tx.QueryRow(`
		INSERT INTO xgstat_fixtures (
			gameweek, fixture_id, fixture_date, 
			home_team, away_team, 
			home_score, away_score, 
			home_xg, away_xg
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (fixture_id, gameweek) 
		DO UPDATE SET
			fixture_date = EXCLUDED.fixture_date,
			home_team = EXCLUDED.home_team,
			away_team = EXCLUDED.away_team,
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			home_xg = EXCLUDED.home_xg,
			away_xg = EXCLUDED.away_xg,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, fixture.Gameweek, fixture.ID, fixture.Date,
		fixture.HomeTeam, fixture.AwayTeam,
		fixture.HomeScore, fixture.AwayScore,
		fixture.HomeXG, fixture.AwayXG,
	).Scan(&fixtureID)

	if err != nil {
		return fmt.Errorf("failed to insert fixture: %w", err)
	}

	// Delete existing shots for this fixture to avoid duplicates
	_, err = tx.Exec("DELETE FROM xgstat_shots WHERE fixture_id = $1", fixtureID)
	if err != nil {
		return fmt.Errorf("failed to delete existing shots: %w", err)
	}

	// Insert home shots
	for _, shot := range fixture.HomeShots {
		err = s.insertShot(tx, fixtureID, shot, "home")
		if err != nil {
			return fmt.Errorf("failed to insert home shot: %w", err)
		}
	}

	// Insert away shots
	for _, shot := range fixture.AwayShots {
		err = s.insertShot(tx, fixtureID, shot, "away")
		if err != nil {
			return fmt.Errorf("failed to insert away shot: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// insertShot inserts a single shot record
func (s *Service) insertShot(tx *sql.Tx, fixtureID int, shot domain.DBXGStatShot, teamType string) error {
	_, err := tx.Exec(`
		INSERT INTO xgstat_shots (
			fixture_id, x, y, xg, is_goal, 
			shot_type, player_name, minute, team_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, fixtureID, shot.X, shot.Y, shot.XG, shot.IsGoal,
		shot.ShotType, shot.PlayerName, shot.Minute, teamType,
	)
	return err
}

// GetFixtureByID retrieves a fixture with its shots by fixture ID
func (s *Service) GetFixtureByID(fixtureID int) (*domain.DBXGStatFixture, error) {
	var fixture domain.DBXGStatFixture
	var dbID int

	err := s.db.QueryRow(`
		SELECT id, gameweek, fixture_id, fixture_date,
			   home_team, away_team, home_score, away_score,
			   home_xg, away_xg
		FROM xgstat_fixtures
		WHERE fixture_id = $1
	`, fixtureID).Scan(
		&dbID, &fixture.Gameweek, &fixture.ID, &fixture.Date,
		&fixture.HomeTeam, &fixture.AwayTeam,
		&fixture.HomeScore, &fixture.AwayScore,
		&fixture.HomeXG, &fixture.AwayXG,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("fixture not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query fixture: %w", err)
	}

	// Get shots
	rows, err := s.db.Query(`
		SELECT x, y, xg, is_goal, shot_type, player_name, minute, team_type
		FROM xgstat_shots
		WHERE fixture_id = $1
		ORDER BY minute
	`, dbID)
	if err != nil {
		return nil, fmt.Errorf("failed to query shots: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var shot domain.DBXGStatShot
		var teamType string
		err := rows.Scan(
			&shot.X, &shot.Y, &shot.XG, &shot.IsGoal,
			&shot.ShotType, &shot.PlayerName, &shot.Minute, &teamType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shot: %w", err)
		}

		if teamType == "home" {
			fixture.HomeShots = append(fixture.HomeShots, shot)
		} else {
			fixture.AwayShots = append(fixture.AwayShots, shot)
		}
	}

	return &fixture, nil
}
