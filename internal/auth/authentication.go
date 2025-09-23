package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// hash-password
func HashPassword(password string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// check-password-hash
func CheckPasswordHash(password, hash string) error{
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// make-jwt
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	now := time.Now().UTC()

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Issuer:    "chirpy",
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// validate-jwt
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error){
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