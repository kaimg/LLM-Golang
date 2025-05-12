package handlers

import (
	"llm/api"
	"llm/config"
	"llm/db"
	"llm/logger"
	"llm/utils"
	"net/http"
	"log/slog"
	"fmt"
	"strings"
)

func PromptHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from session
	session, _ := config.SessionStore.Get(r, config.SessionName)
	userID, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Please log in to use this feature", http.StatusUnauthorized)
		return
	}

	// Get user's encrypted API key from database
	var encryptedApiKey string
	err := db.DB.QueryRow("SELECT groq_api_key FROM users WHERE id = $1", userID).Scan(&encryptedApiKey)
	if err != nil {
		logger.Logger.Error("Failed to get API key", slog.String("error", err.Error()))
		http.Error(w, "Failed to get API key", http.StatusInternalServerError)
		return
	}

	if encryptedApiKey == "" {
		errorMsg := `<div class="alert alert-warning" role="alert">
			Please set your GROQ API key in your <a href="/profile" class="alert-link">profile settings</a> to use this feature.
		</div>`
		w.Write([]byte(errorMsg))
		return
	}

	// Decrypt the API key
	apiKey, err := utils.DecryptAPIKey(encryptedApiKey)
	if err != nil {
		logger.Logger.Error("Failed to decrypt API key", slog.String("error", err.Error()))
		http.Error(w, "Failed to decrypt API key", http.StatusInternalServerError)
		return
	}

	prompt := r.FormValue("prompt")
	model := r.FormValue("model")
	if model == "" {
		model = "llama3-8b-8192" // Default model
	}
	logger.Logger.Info("Received prompt", slog.String("prompt", prompt), slog.String("model", model))

	// Make request to the API using user's API key
	response, err := api.MakeGroqRequest(prompt, apiKey, model)
	if err != nil {
		logger.Logger.Error("API request failed", 
			slog.String("error", err.Error()),
			slog.String("prompt", prompt))

		var errorMsg string
		if strings.Contains(err.Error(), "API returned status code 401") {
			errorMsg = `<div class="alert alert-danger" role="alert">
				Invalid API key. Please check your API key in your <a href="/profile" class="alert-link">profile settings</a>.
			</div>`
		} else {
			errorMsg = `<div class="alert alert-danger" role="alert">
				Failed to get response from API. Please try again later.
			</div>`
		}
		w.Write([]byte(errorMsg))
		return
	}

	// Format the response
	formattedResponse := utils.FormatMarkdown(response)
	logger.Logger.Debug("Formatted API response", slog.String("response", formattedResponse))

	// Return the formatted response
	w.Write([]byte(formattedResponse))

	// Save the prompt and response
	if err := db.SavePrompt(userID, prompt, formattedResponse, model); err != nil {
		logger.Logger.Error("Failed to save prompt and response", 
			slog.String("userID", fmt.Sprintf("%v", userID)), 
			slog.String("error", err.Error()))
		return
	}
	logger.Logger.Info("Prompt and response saved successfully")
}
