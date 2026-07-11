package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"main.go/handlers"
	jwtmiddleware "main.go/jwt_middleware"
)
// corsMiddleware lets a browser-based frontend (served from a different
// origin/port than :8080) call this API. Without it, every fetch() from
// the HTML frontend would be blocked by the browser before it even reaches
// these handlers, and preflight OPTIONS requests would 404.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRouter() *mux.Router{

	
	r:=mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/register",handlers.Register).Methods("POST","OPTIONS")
	r.HandleFunc("/login",handlers.Login).Methods("POST","OPTIONS")
	r.Handle("/profile", jwtmiddleware.Middleware(http.HandlerFunc(handlers.View_profile))).Methods("GET","OPTIONS")
	r.Handle("/complete_profile",jwtmiddleware.Middleware(http.HandlerFunc(handlers.Complete_Profile))).Methods("POST","OPTIONS")
	return r
}