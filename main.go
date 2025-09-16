package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var port = ":8080"

func main(){
	cfg := &apiConfig{}
	mux := http.NewServeMux()	
	mux.Handle("GET /app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("")))))
	mux.HandleFunc("GET /healthz", appHandler)
	mux.HandleFunc("GET /metrics", cfg.hitHandler)
	mux.HandleFunc("POST /reset", cfg.resetHandler)
	
	server := &http.Server{
		Addr: port,
		Handler: mux,
	}
log.Fatal(server.ListenAndServe())
}	

func appHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hitHandler(w http.ResponseWriter, r *http.Request){
	hits := cfg.fileserverHits.Load()
	x := fmt.Sprintf("Hits: %v", hits)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(x))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Swap(0)
	//x := fmt.Sprintf("Hits: %v", hits)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits: 0"))
}