package jwtmiddleware

import (
	"context"
	"log"
	"net/http"
	"os"

	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct{
	Id int `json:"id"`
	jwt.RegisteredClaims
}
func init(){
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}
	
}
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))


		claims:=&Claims{}
		// verify token
		
		token, err := jwt.ParseWithClaims(tokenString,claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("Secret")), nil
		})
		if err != nil || token == nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// store id in the request context
		ctx:=context.WithValue(r.Context(),"id",claims.Id)

		// token valid, call next handler
		if next != nil {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}