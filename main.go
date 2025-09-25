package main

import (
	"log"
	"net/http"

	"github.com/Johnermac/http-server/internal/api"
	_ "github.com/lib/pq"
)


var port = ":8080"

func main(){
	cfg := api.NewAPIConfig()

	mux := http.NewServeMux()	
	
	// APP
	mux.Handle("GET /app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("")))))
	
	// API
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChirpsHandler)

	mux.HandleFunc("POST /api/chirps", cfg.CreateChirpHandler)
	mux.HandleFunc("POST /api/users", cfg.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cfg.LoginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", cfg.RevokeTokenHandler)

	// ADMIN
	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.DeleteAllUsersHandler)
	
	server := &http.Server{
		Addr: port,
		Handler: mux,
	}
log.Fatal(server.ListenAndServe())
}	









