package handlers

import (
	"html/template"
	"llm/config"
	"llm/db"
	"net/http"
	"path/filepath"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	// Load the template
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Retrieve the session
	session, _ := config.SessionStore.Get(r, config.SessionName)
	userID, ok := session.Values["user_id"]
	// Default data structure to send to the template
	data := struct {
		LoggedIn  bool
		Username  string
		AvatarURL string
	}{
		LoggedIn: false,
	}

	// If the user is logged in, fetch their details from the database
	if ok {
		var username, avatarURL string
		err := db.DB.QueryRow("SELECT username, avatar_url FROM users WHERE id = $1", userID).Scan(&username, &avatarURL)
		if err == nil {
			data.LoggedIn = true
			data.Username = username
			data.AvatarURL = avatarURL
		}
	}

	// Execute the template with the user data
	tmpl.Execute(w, data)
}
