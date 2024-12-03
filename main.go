package main

import (
	"llm/auth"
	"llm/config"
	"llm/db"
    "llm/handlers"
    "llm/logger"
    "log/slog"
    "fmt"
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

    http.HandleFunc("/", handlers.FormHandler)
    http.HandleFunc("/api/prompt", handlers.PromptHandler)
	http.HandleFunc("/auth/login", auth.LoginHandler)
    http.HandleFunc("/auth/logout", auth.LogoutHandler)
    http.HandleFunc("/auth/callback", auth.CallbackHandler)

    http.HandleFunc("/login", handlers.LoginPageHandler)
    http.HandleFunc("/loginemail", auth.LoginViaEmailHandler)
    address := fmt.Sprintf(":%s", config.Port)

    logger.Logger.Info("Server is running", slog.String("address", address))

    if err := http.ListenAndServe(address, nil); err != nil {
		logger.Logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}