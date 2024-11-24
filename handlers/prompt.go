package handlers

import (
	"llm/api"
	"llm/config"
	"llm/db"
	"llm/utils"
	"log"
	"net/http"
)

func PromptHandler(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("prompt")
	response, err := api.MakeGroqRequest(prompt)
	if err != nil {
		http.Error(w, "Failed to get response from API", http.StatusInternalServerError)
		log.Println("API request error:", err)
		return
	}

	// Format the response
	formattedResponse := utils.FormatMarkdown(response)

	// Return the formatted response
	w.Write([]byte(formattedResponse))

	// Save the prompt and response
	session, _ := config.SessionStore.Get(r, config.SessionName)
	userID := session.Values["user_id"]
	db.SavePrompt(userID, prompt, formattedResponse)
}
