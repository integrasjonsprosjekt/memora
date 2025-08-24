package status

import (
	"encoding/json"
	"log"
	"memora/internal/config"
	"net/http"
	"time"
)

// Groups all the handlers in the package
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", GetStatus)

	return http.StripPrefix("/status", mux)
}

// Returns status of the service, and the services used by it
func GetStatus(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Version: "v1",
		Uptime:  time.Duration(time.Since(config.StartTime).Seconds()),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Println("Failed to encode response:", err)
	}
}
