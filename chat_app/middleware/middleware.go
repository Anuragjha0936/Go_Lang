package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey="userID"

var jwtSecret =[]byte(GetJwtSecret())

func GetJwtSecret() string{
	s:=os.Getenv("Secret")
	if s!=""{
		return s
	}
	return "dev-only-insecure-secret-change-me"
}

func Middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
 
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Verify the signing method is what we expect (HMAC).
			// Without this check, an attacker could craft a token
			// using a different algorithm (e.g. "none", or RSA where
			// the "signature" is just your own public key) and trick
			// naive verification into accepting it. This is the most
			// common real-world JWT vulnerability.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
 
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
 
		// JWT numeric claims decode as float64 in Go, hence the cast.
		// Converted to int to match models.User.Id / models.Message
		// .SenderId / .ReceiverId, which are all plain int here.
		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			http.Error(w, "user_id claim missing", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)
 
		// Attach userID to the request context and pass control on to
		// the actual controller (e.g. controllers.GetUsers).
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}