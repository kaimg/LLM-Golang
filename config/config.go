package config

import (
    "log"
    "os"

    "github.com/markbates/goth"
    "github.com/markbates/goth/gothic"
    "github.com/markbates/goth/providers/github"
    "github.com/gorilla/sessions"
    "github.com/joho/godotenv"
)

var (
    GitHubClientID     string
    GitHubClientSecret string
    GitHubRedirectURL  string
    GitHubAuthURL      string
    GitHubTokenURL     string
    GitHubUserAPIURL   string
    GitHubAPI          string
    
    SessionStore       *sessions.CookieStore
    SessionName        string

    ApiUrl string
    ApiKey string

    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string

    Port string
)

func LoadConfig() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Port
    Port = os.Getenv("PORT")

    // GitHub OAuth
    GitHubClientID = os.Getenv("GITHUB_CLIENT_ID")
    GitHubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
    GitHubRedirectURL = os.Getenv("GITHUB_REDIRECT_URL")
    GitHubAuthURL = os.Getenv("GITHUB_AUTH_URL")
    GitHubTokenURL = os.Getenv("GITHUB_TOKEN_URL")
    GitHubUserAPIURL = os.Getenv("GITHUB_USER_API_URL")
    GitHubAPI = os.Getenv("GITHUB_API")

    // Session
    SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
    gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
    SessionName = "llm-session"

    // API Configuration
    ApiUrl = os.Getenv("API_URL")
    ApiKey = os.Getenv("API_KEY")

    // Database Configuration
    DBHost = os.Getenv("DB_HOST")
    DBPort = os.Getenv("DB_PORT")
    DBUser = os.Getenv("DB_USER")
    DBPassword = os.Getenv("DB_PASSWORD")
    DBName = os.Getenv("DB_NAME")

    // Initialize GitHub OAuth provider using goth
    goth.UseProviders(
        github.New(GitHubClientID, GitHubClientSecret, GitHubRedirectURL),
    )
}
