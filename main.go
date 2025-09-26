package main

import (
	"log"
	"net/http"

	"github.com/Johnermac/http-server/internal/api"
	_ "github.com/lib/pq"
)

var port = ":8080"

func main() {
	cfg := api.NewAPIConfig()

	mux := http.NewServeMux()

	// app
	mux.Handle("GET /app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("")))))

	// misc
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)

	// chirps
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChirpsHandler)
	mux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.DeleteChirpHandler)

	// users
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("PUT /api/users", cfg.UpdateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginUserHandler)
	mux.HandleFunc("POST /admin/reset", cfg.DeleteAllUsersHandler)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.UpdatePremiumUserHandler)

	// token
	mux.HandleFunc("POST /api/refresh", cfg.UpdateTokenHandler)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeTokenHandler)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
