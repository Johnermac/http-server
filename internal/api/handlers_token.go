package api

import (
	"net/http"

	"github.com/Johnermac/http-server/internal/auth"
	"github.com/Johnermac/http-server/internal/helpers"
)

// refresh-token
func (cfg *APIConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type responseBody struct {		
		Token		string `json:"token"`		
	}
	
	refreshToken, err := auth.GetBearerToken(r.Header) 
	if err != nil {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}

	// validate the refreshToken in the databse
	token, err := cfg.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}

	tokenString, err := auth.MakeJWT(
		token.UserID, 
		cfg.JWTSecret,
	)

	if err != nil {
		helpers.RespondWithError(w, 500, "Error in Token creation")
		return
	}	

	// Do something with responseBody		
	helpers.RespondWithJSON(w, 200, responseBody{		
		Token: tokenString})	
}

// revoke-refresh-token
func (cfg *APIConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	refreshToken, err := auth.GetBearerToken(r.Header) 
	if err != nil {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}	

	_, err = cfg.DB.UpdateRevokeAt(r.Context(), refreshToken)
	if err != nil {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}	

	helpers.RespondWithJSON(w, 204, "")	
}