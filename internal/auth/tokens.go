package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Makes a JSON Web Token for a user
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("invalid user id")
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "battle-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// Validates a JWT and returns the user id from the subject claim if the token is valid.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Claims struct for storing info parsed with jwt.ParseWithClaims
	claims := &jwt.RegisteredClaims{}

	// Returns the same key type ([]byte) used when the token was signed.
	// An error will be returned if the token is invalid or has expired.
	keyFunc := func(token *jwt.Token) (any, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(tokenSecret), nil
	}

	// Validate the signature of the JWT
	// and extract the claims into a *jwt.Token struct
	_, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return uuid.Nil, err
	}

	// Extract user id from subject claim
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// Extract the bearer token from the authorization header of a request
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	stripedHeader := strings.TrimPrefix(authHeader, "Bearer ")
	return stripedHeader, nil
}
