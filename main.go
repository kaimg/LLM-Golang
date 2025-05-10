package main

import (
	"fmt"
	"llm/config"
	"llm/db"
	"llm/logger"
	"llm/router"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Initialize logger with desired level
	logger.InitLogger(slog.LevelDebug)

	// Load environment variables
	config.LoadConfig()

	// Connect to the database
	if err := db.Connect(); err != nil {
		logger.Logger.Error("Error connecting to the database", "error", err)
		os.Exit(1)
	}

	// Setup routes
	router.SetupRoutes()

	// Start the server
	address := fmt.Sprintf(":%s", config.Port)
	logger.Logger.Info("Server is running", slog.String("address", address))

	if err := http.ListenAndServe(address, nil); err != nil {
		logger.Logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}