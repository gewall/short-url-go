package pkg

import (
	"crypto/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userId uuid.UUID) (string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	claims := jwt.MapClaims{
		"iss": "auth-services",
		"aud": "client",
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"sub": userId,
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}

func GenerateRefreshToken() (string, error) {
	// refreshToken := make([]byte, 32)
	// rand.Text()
	refreshToken := rand.Text()

	return string(refreshToken), nil
}
