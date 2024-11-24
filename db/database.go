package db

import (
	"database/sql"
	"fmt"
	"llm/config"
	"log"

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
        return fmt.Errorf("failed to connect to the database: %v", err)
    }

    // Test the connection
    if err := DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping the database: %v", err)
    }
	
	rows, err := DB.Query("SELECT * FROM playing_with_neon")
    if err != nil {
        log.Fatalf("Failed to retrieve query: %v", err)
    }
    defer rows.Close()

    // Print rows
    for rows.Next() {
        var col1 int // Replace with your actual column types
        var col2 string    // Replace with your actual column types
		var col3 float64
        err := rows.Scan(&col1, &col2, &col3) // Adjust based on your column count and types
        if err != nil {
            log.Fatalf("Failed to scan row: %v", err)
        }
        fmt.Printf("Row: col1=%d, col2=%s, col3=%f\n", col1, col2, col3)
    }

    if err := rows.Err(); err != nil {
        log.Fatalf("Error during row iteration: %v", err)
    }
    return nil
}
