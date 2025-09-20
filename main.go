package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Johnermac/http-server/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

var port = ":8080"

func main(){
	cfg := NewAPIConfig()

	mux := http.NewServeMux()	
	
	// APP
	mux.Handle("GET /app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("")))))
	
	// API
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getAllChirpsHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpsHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)

	// ADMIN
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.deleteAllUsersHandler)
	
	server := &http.Server{
		Addr: port,
		Handler: mux,
	}
log.Fatal(server.ListenAndServe())
}	

func newDB() *database.Queries {
    godotenv.Load()

    dbURL := os.Getenv("DB_URL")			

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }

    return database.New(db)
}

func NewAPIConfig() *apiConfig {
	return &apiConfig{
		db: newDB(),	
		platform: os.Getenv("PLATFORM"),	
	}
}


func healthHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct {
		Data string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Data string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	// Normal check
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 500, "Couldn't unmarshal parameters")
		return
	}

	// Do something with requestBody

	if len(params.Data) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	new := badWordReplacement(params.Data)	

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
    Body:   new,
    UserID: params.User_id, // UUID from users table
	})

	if err != nil {
		respondWithError(w, 500, "Create chirp error")
		return
	}


	// Do something with requestBody		
	respondWithJSON(w, 201, responseBody{
		Id: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Data: chirp.Body,
		User_id: chirp.UserID,})				
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request){
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Get All Chirps error")
		return
	}

	// respond with array	
	respondWithJSON(w, 200, chirps)				
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct {
		Email string `json:"email"`			
	}
	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	// Normal error handling
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 500, "Couldn't unmarshal parameters")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "Create user error")
		return
	}

	//fmt.Println("User: %v has been created in DB", user)

	// Do something with requestBody		
	respondWithJSON(w, 201, responseBody{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email})	
}


func badWordReplacement (payload string) string{	
  original := strings.Split(payload, " ")  
	out := make([]string, 0, len(original))
	wordsToFilter := []string{"kerfuffle", "sharbert", "fornax"}

	for _, o := range original {
		if slices.Contains(wordsToFilter, strings.ToLower(o)){
			out = append(out, "****")	
		}	else {
			out = append(out, o)	
		}			
	}

	return strings.Join(out, " ")
}


// JSON Helpers

func respondWithJSON(w http.ResponseWriter, code int, payload any) error {
  response, err := json.Marshal(payload)
  if err != nil {
      return err
  }

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
    return respondWithJSON(w, code, map[string]string{"error": msg})
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request){
	hits := cfg.fileserverHits.Load()
	x := fmt.Sprintf(`<html><body><h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p></body></html>`, hits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(x))
}

func (cfg *apiConfig) deleteAllUsersHandler(w http.ResponseWriter, r *http.Request){
	// reset counting	
	cfg.fileserverHits.Swap(0)

	// reset User in DB		
	if cfg.platform == "dev" {
		err := cfg.db.DeleteAllUsers(r.Context())
		if err != nil {
			respondWithError(w, 500 , "Error Deleting Users")
			return			
		}
		respondWithJSON(w, 200, "All Users Deleted")	
		return

	} else {
		respondWithError(w, 403, "Can only delete in DEV environment!")
		return
	}
}