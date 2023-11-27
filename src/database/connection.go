package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// GetConfigFromArgsOrEnv creates a PostgreSQL connection configuration based on command-line arguments or environment variables.
func GetConfigFromArgsOrEnv(args []string) *Config {
	var host, database, user, password string

	if len(args) >= 4 {
		// Use command-line arguments if provided
		host = os.Getenv(args[0])
		database = os.Getenv(args[1])
		user = os.Getenv(args[2])
		password = os.Getenv(args[3])
	}

	return &Config{
		Host:     host,
		Database: database,
		User:     user,
		Password: password,
	}
}

// Config holds PostgreSQL connection configuration
type Config struct {
	Host     string
	Database string
	User     string
	Password string
}

// NewDBPool creates a new PostgreSQL connection pool
func NewDBPool(cfg *Config) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", cfg.Host, cfg.Database, cfg.User, cfg.Password)
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
