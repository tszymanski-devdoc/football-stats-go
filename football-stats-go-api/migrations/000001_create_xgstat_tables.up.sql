-- Create xgstat_fixtures table
CREATE TABLE IF NOT EXISTS xgstat_fixtures (
    id SERIAL PRIMARY KEY,
    gameweek INT NOT NULL,
    fixture_id INT NOT NULL,
    fixture_date TIMESTAMP NOT NULL,
    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,
    home_score INT NOT NULL,
    away_score INT NOT NULL,
    home_xg DECIMAL(5, 2) NOT NULL,
    away_xg DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(fixture_id, gameweek)
);

-- Create index on gameweek for faster queries
CREATE INDEX idx_xgstat_fixtures_gameweek ON xgstat_fixtures(gameweek);

-- Create index on fixture_date for faster queries
CREATE INDEX idx_xgstat_fixtures_date ON xgstat_fixtures(fixture_date);

-- Create xgstat_shots table
CREATE TABLE IF NOT EXISTS xgstat_shots (
    id SERIAL PRIMARY KEY,
    fixture_id INT NOT NULL REFERENCES xgstat_fixtures(id) ON DELETE CASCADE,
    x DECIMAL(6, 3) NOT NULL,
    y DECIMAL(6, 3) NOT NULL,
    xg DECIMAL(5, 3) NOT NULL,
    is_goal BOOLEAN NOT NULL,
    shot_type VARCHAR(100),
    player_name VARCHAR(255),
    minute INT,
    team_type VARCHAR(10) NOT NULL CHECK (team_type IN ('home', 'away')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_coordinates CHECK (x >= 0 AND x <= 100 AND y >= 0 AND y <= 100),
    CONSTRAINT valid_xg CHECK (xg >= 0 AND xg <= 1),
    CONSTRAINT valid_minute CHECK (minute > 0 AND minute <= 120)
);

-- Create index on fixture_id for faster joins
CREATE INDEX idx_xgstat_shots_fixture_id ON xgstat_shots(fixture_id);

-- Create index on player_name for player-specific queries
CREATE INDEX idx_xgstat_shots_player ON xgstat_shots(player_name);

-- Create index on is_goal for filtering goals
CREATE INDEX idx_xgstat_shots_is_goal ON xgstat_shots(is_goal);
