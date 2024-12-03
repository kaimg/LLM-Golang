package auth

import (
	"fmt"
	"llm/config"
	"llm/db"
    "llm/logger"
	"net/http"
	"encoding/json"
    "log/slog"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// LoginHandler redirects the user to GitHub's OAuth2 login page using Goth
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect the user to GitHub login using Goth
	gothic.BeginAuthHandler(w, r)
}

// CallbackHandler handles the GitHub OAuth2 callback, stores user info in the session
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Get user info from the callback request
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Error completing user authentication", http.StatusInternalServerError)
		return
	}
    // Fetch user's primary email from GitHub API
	userEmail, err := getPrimaryEmail(config.GitHubAPI)
	if err != nil {
		http.Error(w, "Error fetching user email", http.StatusInternalServerError)
		return
	}
	// Insert or update the user in the database
	userID, err := upsertUser(user, userEmail)
	if err != nil {
		http.Error(w, "Error saving user info", http.StatusInternalServerError)
		return
	}

	// Store the user ID in the session
	session, _ := config.SessionStore.Get(r, config.SessionName)
	session.Values["user_id"] = userID
	session.Save(r, w)

	// Redirect the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// LogoutHandler handles user logout by deleting the session and redirecting to login page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Logout the user using Goth
    gothic.Logout(w, r)
	// Retrieve the session
	session, _ := config.SessionStore.Get(r, config.SessionName)

    // Optionally, clear the session cookie to fully remove the session from the client
	session.Options.MaxAge = -1 // This will expire the session cookie
	session.Save(r, w)

	// Delete the "user_id" from the session
	delete(session.Values, "user_id")

	// Redirect to the login page (GitHub login flow)
	http.Redirect(w, r, "/", http.StatusFound)
}

// getPrimaryEmail fetches the primary email of the authenticated user from GitHub API
func getPrimaryEmail(accessToken string) (string, error) {
	// GitHub API endpoint to fetch email
	emailAPIURL := "https://api.github.com/user/emails"
	req, err := http.NewRequest("GET", emailAPIURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response to get email
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	err = json.NewDecoder(resp.Body).Decode(&emails)
	if err != nil {
		return "", err
	}

	// Find the primary email
	var userEmail string
	for _, email := range emails {
		if email.Primary {
			userEmail = email.Email
			break
		}
	}

	if userEmail == "" {
		return "", fmt.Errorf("primary email not found")
	}

	return userEmail, nil
}
// upsertUser inserts or updates the user in the database and returns the user ID
func upsertUser(user goth.User, email string) (int, error) {
	var userID int

	// Check if the user already exists
	err := db.DB.QueryRow("SELECT id FROM users WHERE github_id = $1", user.UserID).Scan(&userID)
	if err != nil {
		// Insert new user
		err = db.DB.QueryRow(
			"INSERT INTO users (username, github_id, email, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id",
			user.Name, user.UserID, email, user.AvatarURL,
		).Scan(&userID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert user: %v", err)
		}
	} else {
		// Update existing user
		_, err = db.DB.Exec(
			"UPDATE users SET username = $1, email = $2, avatar_url = $3 WHERE github_id = $4",
			user.Name, email, user.AvatarURL, user.UserID,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to update user: %v", err)
		}
	}

	return userID, nil
}
func LoginViaEmailHandler(w http.ResponseWriter, r *http.Request) {
	var userID = 1
	var email = r.FormValue("email")
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
        logger.Logger.Debug("DB", slog.String("userID", fmt.Sprintf("%v", err)))
        logger.Logger.Debug("Email", slog.String("email", fmt.Sprintf("%v", email)))
		//http.Error(w, "Error while reading DB for users", http.StatusInternalServerError)
		return
	}
	if userID != 0 {
		logger.Logger.Debug("Email Login", slog.String("userID", fmt.Sprintf("%v", userID)))
		
	}
    // Retrieve the session
	session, _ := config.SessionStore.Get(r, config.SessionName)
    session.Values["user_id"] = userID
	session.Save(r, w)

	// Redirect the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

