package main

import (
	"html/template"
	"llm/api"
	"llm/auth"
	"llm/config"
	"llm/db"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

)

func main() {
    // Load environment variables
    config.LoadConfig()
    
    // Connect to the database
    if err := db.Connect(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    http.HandleFunc("/", formHandler)
    http.HandleFunc("/api/prompt", promptHandler)
	http.HandleFunc("/auth/login", auth.LoginHandler)
    http.HandleFunc("/auth/callback", auth.CallbackHandler)

    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    tmplPath := filepath.Join("templates", "index.html")
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }

    session, _ := config.SessionStore.Get(r, config.SessionName)
    userID, ok := session.Values["user_id"]

    data := struct {
        LoggedIn  bool
        Username  string
        AvatarURL string
    }{
        LoggedIn: false,
    }

    if ok {
        var username, avatarURL string
        err := db.DB.QueryRow("SELECT username, avatar_url FROM users WHERE id = $1", userID).Scan(&username, &avatarURL)
        if err == nil {
            data.LoggedIn = true
            data.Username = username
            data.AvatarURL = avatarURL
        }
    }

    tmpl.Execute(w, data)
}

func promptHandler(w http.ResponseWriter, r *http.Request) {
    prompt := r.FormValue("prompt")
    response, err := api.MakeGroqRequest(prompt)
    if err != nil {
        http.Error(w, "Failed to get response from API", http.StatusInternalServerError)
        log.Println("API request error:", err)
        return
    }

    // Format the response to HTML
    formattedResponse := formatMarkdown(response)

    // Return the formatted response
    w.Write([]byte(formattedResponse))
}

// formatMarkdown converts a Markdown-style response to basic HTML.
func formatMarkdown(text string) string {
    // Convert headers (e.g., # Header) to <h1>, <h2>, etc.
    text = regexp.MustCompile(`(?m)^### (.+)$`).ReplaceAllString(text, "<h3>$1</h3>")
    text = regexp.MustCompile(`(?m)^## (.+)$`).ReplaceAllString(text, "<h2>$1</h2>")
    text = regexp.MustCompile(`(?m)^# (.+)$`).ReplaceAllString(text, "<h1>$1</h1>")

    // Convert bold (**text**) to <b>text</b>
    text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "<b>$1</b>")

    // Convert italic (*text*) to <i>text</i>
    text = regexp.MustCompile(`\*(.+?)\*`).ReplaceAllString(text, "<i>$1</i>")

	// Convert code (```text```) to <code>text<code>
	text = regexp.MustCompile("(?s)```(.*?)```").ReplaceAllString(text, "<code>$1</code>")
	
    // Convert newline characters to <br> for line breaks
    text = strings.ReplaceAll(text, "\n", "<br>")

    return text
}
