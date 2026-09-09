package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"main.go/extractor"
	"main.go/parser"
)


func writeJSON(w http.ResponseWriter, status int, payload map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "failed to parse form"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "no file uploaded"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": "failed to read file"})
		return
	}

	text, err := extractor.Extract(header.Filename, data)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	if strings.TrimSpace(text) == "" {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{"error": "no text found in file"})
		return
	}

	resume, err := parser.ParseResume(text)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "success", "data": resume})
}

func main() {
	// ── API connection check on startup ──
	fmt.Println("Connecting to NVIDIA API...")
	if err := parser.CheckConnection(); err != nil {
		fmt.Println("Connection failed:", err)
		// don't exit — still start server, just warn
	} else {
		fmt.Println("Connected to NVIDIA API successfully!")
	}

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/parse", parseHandler).Methods(http.MethodPost, http.MethodOptions)

	fmt.Println("Listening on :8080")
	
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Server failed:", err)
	}
}
