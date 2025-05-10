package db

import (
	"database/sql"
	"fmt"
	"llm/config"
	"llm/logger"
    "log/slog"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect initializes the database connection
func Connect() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		logger.Logger.Error("Failed to connect to the database", slog.String("error", err.Error()))
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		logger.Logger.Error("Failed to ping the database", slog.String("error", err.Error()))
		return fmt.Errorf("failed to ping the database: %v", err)
	}
	logger.Logger.Info("Database connected successfully")

	return nil
}
