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

func init(){
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}
}
func main()  {
	hub:=WS.NewHub()
	r:=router.SetUpRoute(hub)
	fmt.Println("Server is Running")
	port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}

log.Fatal(http.ListenAndServe(":"+port, r))
}
