package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	App       AppConfig
	Server    ServerConfig
	Database  DatabaseConfig
	Ingestion IngestionConfig
	Analytics AnalyticsConfig
}

// AppConfig holds application-level settings
type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds DB connection settings
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// IngestionConfig holds football data ingestion settings
type IngestionConfig struct {
	Provider       string
	APIKey         string
	BaseURL        string
	PollInterval   time.Duration
	RequestTimeout time.Duration
}

// AnalyticsConfig holds analysis-specific tuning parameters
type AnalyticsConfig struct {
	XGModel string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "football-analytics"),
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "0.1.0"),
		},
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: getDuration("SERVER_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "postgres"),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "football"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		Ingestion: IngestionConfig{
			Provider:       getEnv("INGESTION_PROVIDER", "football-data"),
			APIKey:         getEnv("INGESTION_API_KEY", ""),
			BaseURL:        getEnv("INGESTION_BASE_URL", ""),
			PollInterval:   getDuration("INGESTION_POLL_INTERVAL", 10*time.Minute),
			RequestTimeout: getDuration("INGESTION_REQUEST_TIMEOUT", 10*time.Second),
		},
		Analytics: AnalyticsConfig{
			XGModel: getEnv("ANALYTICS_XG_MODEL", "default"),
		},
	}

	validate(cfg)
	return cfg
}

// --- helpers ---

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// validate ensures required configuration is present
func validate(cfg *Config) {
	if cfg.Database.Host == "" || cfg.Database.Name == "" {
		log.Fatal("database configuration is invalid")
	}

	if cfg.Server.Port == "" {
		log.Fatal("server port must be set")
	}
}
