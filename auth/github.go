package auth

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "llm/config"
    "llm/db"
    "net/http"
    "net/url"
)

// LoginHandler redirects the user to GitHub's OAuth2 login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email",
        config.GitHubAuthURL,
        config.GitHubClientID,
        url.QueryEscape(config.GitHubRedirectURL),
    )
    http.Redirect(w, r, authURL, http.StatusFound)
}

// CallbackHandler handles the GitHub OAuth2 callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Code not provided", http.StatusBadRequest)
        return
    }

    token, err := exchangeCodeForToken(code)
    if err != nil {
        http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
        return
    }

    userInfo, err := getUserInfo(token)
    if err != nil {
        http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
        return
    }

    // Extract user data
    username := userInfo["login"].(string)
    githubID := fmt.Sprintf("%v", userInfo["id"])
    email := ""
    if userInfo["email"] != nil {
        email = userInfo["email"].(string)
    }
    avatarURL := ""
    if userInfo["avatar_url"] != nil {
        avatarURL = userInfo["avatar_url"].(string)
    }

    // Insert or update user in the database
    userID, err := upsertUser(username, githubID, email, avatarURL)
    if err != nil {
        http.Error(w, "Failed to save user", http.StatusInternalServerError)
        return
    }

    // Store the user ID in the session
    session, _ := config.SessionStore.Get(r, config.SessionName)
    session.Values["user_id"] = userID
    session.Save(r, w)

    http.Redirect(w, r, "/", http.StatusFound)
}

func exchangeCodeForToken(code string) (string, error) {
    data := url.Values{}
    data.Set("client_id", config.GitHubClientID)
    data.Set("client_secret", config.GitHubClientSecret)
    data.Set("code", code)
    data.Set("redirect_uri", config.GitHubRedirectURL)

    resp, err := http.PostForm(config.GitHubTokenURL, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    values, err := url.ParseQuery(string(body))
    if err != nil {
        return "", err
    }

    return values.Get("access_token"), nil
}

func getUserInfo(token string) (map[string]interface{}, error) {
    req, err := http.NewRequest("GET", config.GitHubUserAPIURL, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        return nil, err
    }

    return userInfo, nil
}

func upsertUser(username, githubID, email, avatarURL string) (int, error) {
    var userID int

    // Check if the user already exists
    err := db.DB.QueryRow(
        "SELECT id FROM users WHERE github_id = $1",
        githubID,
    ).Scan(&userID)

    if err == sql.ErrNoRows {
        // Insert new user
        err = db.DB.QueryRow(
            "INSERT INTO users (username, github_id, email, avatar_url) VALUES ($1, $2, $3, $4) RETURNING id",
            username, githubID, email, avatarURL,
        ).Scan(&userID)
    } else if err == nil {
        // Update existing user
        _, err = db.DB.Exec(
            "UPDATE users SET username = $1, email = $2, avatar_url = $3 WHERE github_id = $4",
            username, email, avatarURL, githubID,
        )
    }

    if err != nil {
        return 0, fmt.Errorf("failed to upsert user: %v", err)
    }

    return userID, nil
}
