package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"main.go/WS"
	"main.go/router"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}
}

func main()  {
	hub := WS.NewHub()
	r := router.SetUpRoute(hub)

	port := os.Getenv("PORT")

	if port == "" {
		port = "10000"
	}

	fmt.Println("Server is Running on port:", port)

	err := http.ListenAndServe("0.0.0.0:"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
