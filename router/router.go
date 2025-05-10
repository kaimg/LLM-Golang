package router

import (
	"llm/auth"
	"llm/handlers"
	"llm/logger"
	"log/slog"
	"net/http"
)

// Route represents a single route in the application
type Route struct {
	Path        string
	Handler     http.HandlerFunc
	Method      string
	Protected   bool
	Description string
}

// Routes groups routes by their function
type Routes struct {
	Auth     []Route
	API      []Route
	Pages    []Route
	Profile  []Route
}

// InitializeRoutes defines all application routes
func InitializeRoutes() Routes {
	return Routes{
		Auth: []Route{
			{
				Path:        "/auth/login",
				Handler:     auth.LoginHandler,
				Method:      "GET",
				Protected:   false,
				Description: "GitHub OAuth login endpoint",
			},
			{
				Path:        "/auth/logout",
				Handler:     auth.LogoutHandler,
				Method:      "GET",
				Protected:   true,
				Description: "Logout endpoint",
			},
			{
				Path:        "/auth/callback",
				Handler:     auth.CallbackHandler,
				Method:      "GET",
				Protected:   false,
				Description: "GitHub OAuth callback endpoint",
			},
			{
				Path:        "/loginemail",
				Handler:     auth.LoginViaEmailHandler,
				Method:      "POST",
				Protected:   false,
				Description: "Email login endpoint",
			},
		},
		API: []Route{
			{
				Path:        "/api/prompt",
				Handler:     handlers.PromptHandler,
				Method:      "POST",
				Protected:   true,
				Description: "LLM prompt endpoint",
			},
		},
		Pages: []Route{
			{
				Path:        "/",
				Handler:     handlers.FormHandler,
				Method:      "GET",
				Protected:   false,
				Description: "Home page",
			},
			{
				Path:        "/login",
				Handler:     handlers.LoginPageHandler,
				Method:      "GET",
				Protected:   false,
				Description: "Login page",
			},
		},
		Profile: []Route{
			{
				Path:        "/profile",
				Handler:     handlers.ProfilePageHandler,
				Method:      "GET",
				Protected:   true,
				Description: "User profile page",
			},
			{
				Path:        "/profile/update-api-key",
				Handler:     handlers.UpdateApiKeyHandler,
				Method:      "POST",
				Protected:   true,
				Description: "Update GROQ API key endpoint",
			},
		},
	}
}

// applyMiddleware applies the necessary middleware stack to a handler
func applyMiddleware(route Route) http.HandlerFunc {
	handler := route.Handler

	// Apply method restriction
	handler = MethodMiddleware(route.Method, handler)

	// Apply authentication if required
	if route.Protected {
		handler = AuthMiddleware(handler)
	}

	// Apply logging (always last to catch all requests)
	handler = LoggingMiddleware(handler)

	return handler
}

// SetupRoutes registers all application routes
func SetupRoutes() {
	routes := InitializeRoutes()
	logger.Logger.Info("Setting up routes...")

	// Helper function to register routes
	registerRoutes := func(routes []Route, groupName string) {
		for _, route := range routes {
			handler := applyMiddleware(route)
			http.HandleFunc(route.Path, handler)
			logger.Logger.Debug("Registered route",
				slog.String("group", groupName),
				slog.String("path", route.Path),
				slog.String("method", route.Method),
				slog.Bool("protected", route.Protected))
		}
	}

	// Register all route groups
	registerRoutes(routes.Auth, "Auth")
	registerRoutes(routes.API, "API")
	registerRoutes(routes.Pages, "Pages")
	registerRoutes(routes.Profile, "Profile")

	logger.Logger.Info("All routes registered successfully")
} 