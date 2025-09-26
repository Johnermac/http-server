package api

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/Johnermac/http-server/internal/database"
	"github.com/Johnermac/http-server/internal/helpers"
	"github.com/google/uuid"
)

// create-chirp
func (cfg *APIConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct {
		Data string `json:"body"`		
	}
	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Data string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	// Parse request
	params, err := helpers.ParseRequest[requestBody](r)
	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	// Auth
	userID, err := cfg.AuthenticateRequest(r)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}

	// Business logic
	if len(params.Data) > 140 {
		helpers.RespondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
    Body:   helpers.BadWordReplacement(params.Data),
    UserID: userID, // UUID from users table
	})

	if err != nil {
		helpers.RespondWithError(w, 500, "Create chirp error")
		return
	}


	// Do something with responseBody		
	helpers.RespondWithJSON(w, 201, responseBody{
		Id: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Data: chirp.Body,
		User_id: chirp.UserID,})				
}

// get-all-chirps
func (cfg *APIConfig) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request){
	var chirps []database.Chirp
	var err error

	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Data string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	authorID := r.URL.Query().Get("author_id")
	if len(authorID) < 1 {
		// author is not specified, get-all-chirps
		chirps, err = cfg.DB.GetAllChirps(r.Context())
	} else {
		userID,err := uuid.Parse(authorID)
		if err != nil {
			helpers.RespondWithError(w, 500, "Error parsing authorID")
			return
		}	
		// get-chirps-by-author
		chirps, err = cfg.DB.GetChirpsByAuthor(r.Context(), userID)	
	}
	
	if err != nil {
		helpers.RespondWithError(w, 500, "Get Chirps error")
		return
	}

	sortParam := r.URL.Query().Get("sort")
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt) // desc
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) // asc
		})
	}

	responses := make([]responseBody, len(chirps))
	for i, c := range chirps {
		responses[i] = responseBody{
			Id:        c.ID,
			Created_at: c.CreatedAt,
			Updated_at: c.UpdatedAt,
			Data:      c.Body,
			User_id:    c.UserID,
		}
	}

	helpers.RespondWithJSON(w, 200, responses)
					
}

// get-chirp
func (cfg *APIConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request){
	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Data string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	chirpIDStr := r.PathValue("chirpID")

	// Convert string → UUID
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid chirp ID")
		return
	}
	
	chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		helpers.RespondWithError(w, 404, "Chirp Not Found")
		return
	}

	// respond with responseBody	
	helpers.RespondWithJSON(w, 200, responseBody{
		Id: chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Data: chirp.Body,
		User_id: chirp.UserID,
	})				
}

// delete-chirp
func (cfg *APIConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request){
	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Data string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	chirpIDStr := r.PathValue("chirpID")

	// Convert string → UUID
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		helpers.RespondWithError(w, 400, "Invalid chirp ID")
		return
	}	

	// Auth
	userID, err := cfg.AuthenticateRequest(r)
	if err != nil {
		helpers.RespondWithError(w, 401, err.Error())
		return
	}	
	
	chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if errors.Is(err, sql.ErrNoRows) { helpers.RespondWithError(w, 404, "Chirp not found"); return }
	if err != nil {
		helpers.RespondWithError(w, 500, "Database error")
		return
	}

	if chirp.UserID != userID {
		helpers.RespondWithError(w, 403, "Forbidden")
		return
	}

	err = cfg.DB.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userID,
		ID: chirpID,
	})
	if err != nil {
		helpers.RespondWithError(w, 500, "Database error")
		return
	}

	// respond with responseBody	
	helpers.RespondNoContent(w)			
}