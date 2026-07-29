package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/joho/godotenv"
	"main.go/WS"
	"main.go/router"
)

func init(){
	err:=godotenv.Load()
	if err!=nil{
		log.Println(".env file not found, using Railway environment variables")
	}
}
func main()  {
	hub:=WS.NewHub()
	r:=router.SetUpRoute(hub)
	fmt.Println("Server is Running")
	err:=http.ListenAndServe(":8080",r)
	if err!=nil{
		log.Fatal(err)
	}
}
