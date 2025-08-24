package router

import (
	"net/http"

	"memora/internal/config"
	"memora/internal/handlers/status"
	"memora/internal/middleware"
)

// New creates a new HTTP handler for the API
func New() http.Handler {
	mux := http.NewServeMux()

	v1 := http.NewServeMux()
	v1.Handle("/status/", status.Handler())

	middlewares := []func(http.Handler) http.Handler{
		middleware.CORS,
		middleware.JSON,
	}

	// Only add HTTP logging in debug mode
	if config.CurrentLevel == config.LogLevelDebug {
		middlewares = append(middlewares, middleware.Logging)
	}

	// Apply middleware to v1 routes
	v1Handler := middleware.Chain(middlewares...)(v1)

	// Mount v1 routes under /api/v1
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Handler))

	return mux
}
