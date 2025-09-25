package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Johnermac/http-server/internal/auth"
	"github.com/google/uuid"
)

// authenticate-request
func (cfg *APIConfig) AuthenticateRequest(r *http.Request) (uuid.UUID, error) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Missing or invalid Authorization header")
	}

	// sanity check
	if strings.Count(tokenString, ".") != 2 {
		return uuid.Nil, fmt.Errorf("Invalid or expired token")
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.JWTSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid or expired token")
	}
	return userID, nil
}