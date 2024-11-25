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
)

func PromptHandler(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("prompt")
	logger.Logger.Info("Received prompt", slog.String("prompt", prompt))

	// Make request to the API
	response, err := api.MakeGroqRequest(prompt)
	if err != nil {
		http.Error(w, "Failed to get response from API", http.StatusInternalServerError)
		logger.Logger.Error("Failed to get response from API", slog.String("prompt", prompt), slog.String("error", err.Error()))
		return
	}

	// Format the response
	formattedResponse := utils.FormatMarkdown(response)
	logger.Logger.Debug("Formatted API response", slog.String("response", formattedResponse))

	// Return the formatted response
	w.Write([]byte(formattedResponse))

	// Save the prompt and response
	session, _ := config.SessionStore.Get(r, config.SessionName)
	userID := session.Values["user_id"]
	logger.Logger.Info("Saving prompt and response", slog.String("userID", fmt.Sprintf("%v", userID)))

	if err := db.SavePrompt(userID, prompt, formattedResponse); err != nil {
		logger.Logger.Error("Failed to save prompt and response", slog.String("userID", fmt.Sprintf("%v", userID)), slog.String("error", err.Error()))
		return
	}
	logger.Logger.Info("Prompt and response saved successfully")
}
