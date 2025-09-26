package api

import (
	"fmt"
	"net/http"

	"github.com/Johnermac/http-server/internal/helpers"
)

// health-check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, 200, "OK")
}

// metrics-handler
func (cfg *APIConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.FileserverHits.Load()
	x := fmt.Sprintf(`<html><body><h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p></body></html>`, hits)

	helpers.RespondWithJSON(w, 200, x)
}
