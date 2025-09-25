package api

import (
	"net/http"
	"time"

	"github.com/Johnermac/http-server/internal/auth"
	"github.com/Johnermac/http-server/internal/database"
	"github.com/Johnermac/http-server/internal/helpers"
	"github.com/google/uuid"
)

// create-user
func (cfg *APIConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct {		
		Email string `json:"email"`			
		Password string `json:"password"`
	}
	type responseBody struct {		
		Id uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	// Parse request
	params, err := helpers.ParseRequest[requestBody](r)
	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		helpers.RespondWithError(w, 500, "Error with Hash Password")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		helpers.RespondWithError(w, 500, "Create user error")
		return
	}

	//fmt.Println("User: %v has been created in DB", user)

	// Do something with requestBody		
	helpers.RespondWithJSON(w, 201, responseBody{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email})	
}

// login-user
func (cfg *APIConfig) LoginHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct {		
		Email			string `json:"email"`			
		Password 	string `json:"password"`
		// Expires_in_seconds  int `json:"expires_in_seconds "`
	}
	type responseBody struct {		
		Id 						uuid.UUID `json:"id"`
		Created_at 		time.Time `json:"created_at"`
		Updated_at 		time.Time `json:"updated_at"`
		Email 				string `json:"email"`
		Token 				string `json:"token"`
		Refresh_token string `json:"refresh_token"`
		
	}

	// Parse request
	params, err := helpers.ParseRequest[requestBody](r)
	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}	

	user, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		helpers.RespondWithError(w, 404, "User not Found")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}			

	tokenString, err := auth.MakeJWT(
		user.ID, 
		cfg.JWTSecret,
	)

	if err != nil {
		helpers.RespondWithError(w, 500, "Error in Token creation")
		return
	}

	// get refresh token
	refreshToken, err := auth.MakeRefreshToken() 
	if err != nil {
		helpers.RespondWithError(w, 500, "Error in Refresh Token creation")
		return
	}

	// save-refresh-token in DB
	cfg.DB.InsertRefreshToken(r.Context(), database.InsertRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	})


	// Do something with requestBody		
	helpers.RespondWithJSON(w, 200, responseBody{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
		Token: tokenString,
		Refresh_token: refreshToken,})	
}

// delete-all-users
func (cfg *APIConfig) DeleteAllUsersHandler(w http.ResponseWriter, r *http.Request){
	// reset counting	
	cfg.FileserverHits.Swap(0)

	// reset User in DB		
	if cfg.Platform == "dev" {
		err := cfg.DB.DeleteAllUsers(r.Context())
		if err != nil {
			helpers.RespondWithError(w, 500 , "Error Deleting Users")
			return			
		}
		helpers.RespondWithJSON(w, 200, "All Users Deleted")	
		return

	} else {
		helpers.RespondWithError(w, 403, "Can only delete in DEV environment!")
		return
	}
}