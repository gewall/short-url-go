package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

		tokenRaw := r.Header.Get("Authorization")
		if tokenRaw == "" {
			render.Render(w, r, pkg.Unauthorized(pkg.ErrUnauthorized))
			return
		}
		tokenParts := strings.Split(tokenRaw, " ")

		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			render.Render(w, r, pkg.Unauthorized(pkg.ErrInvalidToken))
			return
		}
		tokenUser := tokenParts[1]
		token, err := jwt.Parse(tokenUser, func(token *jwt.Token) (any, error) {
			return []byte(jwtKey), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			render.Render(w, r, pkg.Unauthorized(err))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			render.Render(w, r, pkg.Unauthorized(pkg.ErrInvalidToken))
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims["sub"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
