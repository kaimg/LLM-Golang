package handlers

import (
	"fmt"
	"html/template"
	"llm/config"
	"llm/db"
	"llm/utils"
	"llm/logger"
	"log/slog"
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
		DefaultModel string
	}{
		LoggedIn: false,
	}

	// If the user is logged in, fetch their details from the database
	if ok {
		var username, avatarURL, defaultModel string
		err := db.DB.QueryRow("SELECT username, avatar_url, default_model FROM users WHERE id = $1", userID).Scan(&username, &avatarURL, &defaultModel)
		if err == nil {
			data.LoggedIn = true
			data.Username = username
			data.AvatarURL = avatarURL
			data.DefaultModel = defaultModel
		}
	}
	// Execute the template with the user data
	tmpl.Execute(w, data)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	// Load the template
	tmplPath := filepath.Join("templates", "login.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
    data := 0
    tmpl.Execute(w, data)
}

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	// Load the template
	tmplPath := filepath.Join("templates", "profile.html")
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
		LoggedIn    bool
		Username    string
		AvatarURL   string
		Email       string
		CreatedAt   string
		GroqApiKey  string
		DefaultModel string
	}{
		LoggedIn: false,
	}

	// If the user is logged in, fetch their details from the database
	if ok {
		var username, avatarURL, createdAt, email, encryptedApiKey, defaultModel string
		err := db.DB.QueryRow("SELECT username, avatar_url, created_at, email, groq_api_key, default_model FROM users WHERE id = $1", userID).Scan(&username, &avatarURL, &createdAt, &email, &encryptedApiKey, &defaultModel)
		createdAt = utils.ParseDateAndTime(createdAt)
		if err == nil {
			data.LoggedIn = true
			data.Username = username
			data.AvatarURL = avatarURL
			data.Email = email
			data.CreatedAt = createdAt
			data.DefaultModel = defaultModel
			
			// Decrypt the API key if it exists
			if encryptedApiKey != "" {
				if apiKey, err := utils.DecryptAPIKey(encryptedApiKey); err == nil {
					data.GroqApiKey = apiKey
					logger.Logger.Debug("API key decrypted successfully", 
						slog.String("userID", fmt.Sprintf("%v", userID)),
						slog.String("keyLength", fmt.Sprintf("%d", len(apiKey))))
				} else {
					logger.Logger.Error("Failed to decrypt API key", slog.String("error", err.Error()))
				}
			}
		}
	}
	tmpl.Execute(w, data)
}

func UpdateApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from session
	session, _ := config.SessionStore.Get(r, config.SessionName)
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the API key from the form
	groqApiKey := r.FormValue("groq_api_key")
	
	// Encrypt the API key
	encryptedKey, err := utils.EncryptAPIKey(groqApiKey)
	if err != nil {
		logger.Logger.Error("Failed to encrypt API key", slog.String("error", err.Error()))
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	logger.Logger.Debug("Updating API key", 
		slog.String("userID", fmt.Sprintf("%v", userID)),
		slog.String("keyLength", fmt.Sprintf("%d", len(groqApiKey))))

	// Update the encrypted API key in the database
	_, err = db.DB.Exec("UPDATE users SET groq_api_key = $1 WHERE id = $2", encryptedKey, userID)
	if err != nil {
		logger.Logger.Error("Failed to update API key in database", slog.String("error", err.Error()))
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("API key updated successfully", slog.String("userID", fmt.Sprintf("%v", userID)))

	// Redirect back to the profile page
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func UpdateDefaultModelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from session
	session, _ := config.SessionStore.Get(r, config.SessionName)
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the default model from the form
	defaultModel := r.FormValue("default_model")
	
	logger.Logger.Debug("Updating default model", 
		slog.String("userID", fmt.Sprintf("%v", userID)),
		slog.String("model", defaultModel))

	// Update the default model in the database
	_, err := db.DB.Exec("UPDATE users SET default_model = $1 WHERE id = $2", defaultModel, userID)
	if err != nil {
		logger.Logger.Error("Failed to update default model in database", slog.String("error", err.Error()))
		http.Error(w, "Failed to update default model", http.StatusInternalServerError)
		return
	}

	logger.Logger.Info("Default model updated successfully", slog.String("userID", fmt.Sprintf("%v", userID)))

	// Redirect back to the profile page
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}