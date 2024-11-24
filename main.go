package main

import (
	"llm/auth"
	"llm/config"
	"llm/db"
    "llm/handlers"
    "fmt"
	"log"
	"net/http"

)

func main() {
    // Load environment variables
    config.LoadConfig()
    
    // Connect to the database
    if err := db.Connect(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    http.HandleFunc("/", handlers.FormHandler)
    http.HandleFunc("/api/prompt", handlers.PromptHandler)
	http.HandleFunc("/auth/login", auth.LoginHandler)
    http.HandleFunc("/auth/callback", auth.CallbackHandler)

    address := fmt.Sprintf(":%s", config.Port)
    log.Printf("Server started at http://localhost%s", address)
    log.Fatal(http.ListenAndServe(address, nil))
}