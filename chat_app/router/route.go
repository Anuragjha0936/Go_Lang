package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"main.go/WS"
	"main.go/controllers"
	"main.go/middleware"
)

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

func SetUpRoute(hub *WS.Hub) *mux.Router{
	
	// main router
	mr:=mux.NewRouter()

	mr.Use(corsMiddleware)
	// public routes
	mr.HandleFunc("/api/register",controllers.Register).Methods("POST","OPTIONS")
	mr.HandleFunc("/api/login",controllers.Login).Methods("POST","OPTIONS")

	mr.HandleFunc("/ws",WS.ServeWS(hub))

	api:=mr.PathPrefix("/api").Subrouter()
	api.Use(middleware.Middleware)
	api.HandleFunc("/users",controllers.GetUser).Methods("GET","OPTIONS")
	api.HandleFunc("/messages/{userId}",controllers.GetMessage).Methods("GET","OPTIONS")
	return mr
}