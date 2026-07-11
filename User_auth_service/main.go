package main

import (
	"fmt"
	"log"
	"net/http"
	"main.go/routers"
)
func main(){
	r:=routers.SetupRouter()
	fmt.Println("Server is Running")
	err:=http.ListenAndServe(":8080",r)
	if err!=nil{
		log.Fatal(err)
	}
}