package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// hash-password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// check-password-hash
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// make-jwt
func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	now := time.Now().UTC()

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// validate-jwt
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("Invalid subject claim")
	}

	return userId, nil
}

// get-bearer-token
func GetBearerToken(headers http.Header) (string, error) {
	// make sure is not empty
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Empty auth header")
	}

	// make sure its bearer
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("Not a bearer token")
	}

	// make sure the second part is not empty
	if strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("Empty bearer token")
	}

	return parts[1], nil
}

// make-refresh-token
func MakeRefreshToken() (string, error) {
	key := make([]byte, 32) // 256-bit
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return hex.EncodeToString(key), nil
}

// get-api-key
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Empty auth header")
	}

	// make sure its ApiKey
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "ApiKey") {
		return "", errors.New("Not an apiKey token")
	}

	// make sure the second part is not empty
	if strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("Empty apiKey token")
	}

	return parts[1], nil
}
