package api

import (
	"net/http"

	"github.com/Johnermac/http-server/internal/helpers"
)

// health-check
func HealthHandler(w http.ResponseWriter, r *http.Request){
	helpers.RespondWithJSON(w, 200, "OK")
}