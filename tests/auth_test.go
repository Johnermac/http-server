package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Johnermac/http-server/internal/auth"
)

func TestValidateJWT(t *testing.T) {
	secret := "supersecret"
	userID := uuid.New()

	// helper to create tokens
	makeToken := func(expiration time.Duration, signingSecret string) string {
		claims := &jwt.RegisteredClaims{
			Issuer:    "chirpy",
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiration)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, _ := token.SignedString([]byte(signingSecret))
		return ss
	}

	tests := []struct {
		name       string
		token      string
		expectErr  bool
		expectUser uuid.UUID
	}{
		{
			name:       "valid token",
			token:      makeToken(time.Minute, secret),
			expectErr:  false,
			expectUser: userID,
		},
		{
			name:      "wrong secret",
			token:     makeToken(time.Minute, "wrongsecret"),
			expectErr: true,
		},
		{
			name:      "expired token",
			token:     makeToken(-time.Minute, secret),
			expectErr: true,
		},
		{
			name:      "malformed token",
			token:     "not.a.jwt",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotID, err := auth.ValidateJWT(tc.token, secret)

			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if gotID != tc.expectUser {
				t.Errorf("expected userID %v, got %v", tc.expectUser, gotID)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	secret := "test-secret"
	userID := uuid.New()

	type testCase struct {
		name       string
		userID     uuid.UUID
		secret     string
		expiresIn  time.Duration
		shouldPass bool
	}

	testCases := []testCase{
		{"valid token", userID, secret, time.Hour, true},
		{"expired token", userID, secret, -time.Minute, false},
		{"wrong secret", userID, "wrong-secret", time.Hour, false},
	}

	passCount := 0
	failCount := 0

	for _, tc := range testCases {
		tokenString, err := auth.MakeJWT(tc.userID, tc.secret)
		if err != nil {
			t.Errorf("MakeJWT error: %v", err)
			continue
		}

		// Validate the token
		parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil // always validate with correct secret here
		})

		valid := err == nil && parsedToken.Valid

		if valid != tc.shouldPass {
			failCount++
			t.Errorf(`---------------------------------
Case:       %s
Expecting:  %v
Actual:     %v
Fail
`, tc.name, tc.shouldPass, valid)
		} else {
			passCount++
			fmt.Printf(`---------------------------------
Case:       %s
Expecting:  %v
Actual:     %v
Pass
`, tc.name, tc.shouldPass, valid)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func TestGetBearerToken(t *testing.T) {
	type testCase struct {
		name       string
		header     http.Header
		shouldPass bool
	}

	testCases := []testCase{
		{
			name:       "valid token",
			header:     http.Header{"Authorization": []string{"Bearer token_example"}},
			shouldPass: true,
		},
		{
			name:       "not bearer",
			header:     http.Header{"Authorization": []string{"Basic token_example"}},
			shouldPass: false,
		},
		{
			name:       "empty token",
			header:     http.Header{"Authorization": []string{"Bearer "}},
			shouldPass: false,
		},
		{
			name:       "missing header",
			header:     http.Header{},
			shouldPass: false,
		},
	}

	passCount := 0
	failCount := 0

	for _, tc := range testCases {
		_, err := auth.GetBearerToken(tc.header)
		valid := (err == nil) // if no error, consider it valid

		if valid != tc.shouldPass {
			failCount++
			t.Errorf(`---------------------------------
Case:       %s
Expecting:  %v
Actual:     %v
Fail
`, tc.name, tc.shouldPass, valid)
		} else {
			passCount++
			fmt.Printf(`---------------------------------
Case:       %s
Expecting:  %v
Actual:     %v
Pass
`, tc.name, tc.shouldPass, valid)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}
