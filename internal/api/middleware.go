package api

import (
	"fmt"
	"net/http"

	"github.com/Johnermac/http-server/internal/helpers"
)

// midleware-metrics-inc
func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			cfg.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
	})
}

// metrics-handler
func (cfg *APIConfig) MetricsHandler(w http.ResponseWriter, r *http.Request){
	hits := cfg.FileserverHits.Load()
	x := fmt.Sprintf(`<html><body><h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p></body></html>`, hits)

	helpers.RespondWithJSON(w, 200, x)			
}