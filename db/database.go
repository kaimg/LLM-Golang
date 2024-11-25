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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
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

	rows, err := DB.Query("SELECT * FROM playing_with_neon")
	if err != nil {
		logger.Logger.Error("Failed to retrieve query", slog.String("error", err.Error()))
		return fmt.Errorf("failed to retrieve query: %v", err)
	}
	defer rows.Close()

	// Print rows
	for rows.Next() {
		var col1 int    // Replace with your actual column types
		var col2 string // Replace with your actual column types
		var col3 float64
		err := rows.Scan(&col1, &col2, &col3) // Adjust based on your column count and types
		if err != nil {
			logger.Logger.Error("Failed to scan row", slog.String("error", err.Error()))
			return fmt.Errorf("failed to scan row: %v", err)
		}
		logger.Logger.Debug("Row fetched", slog.Int("col1", col1), slog.String("col2", col2), slog.Float64("col3", col3))
	}

	if err := rows.Err(); err != nil {
		logger.Logger.Error("Error during row iteration", slog.String("error", err.Error()))
		return fmt.Errorf("error during row iteration: %v", err)
	}

	return nil
}
