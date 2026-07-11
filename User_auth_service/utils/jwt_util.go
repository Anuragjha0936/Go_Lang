package utils

import (
	
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	

)


func Generate_token(id int) (string,error){
	claims:=jwt.MapClaims{
		"id":id,
		"exp":time.Now().Add(24*time.Hour).Unix(),
	}
	
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	Secret:=[]byte (os.Getenv("Secret"))
	return token.SignedString(Secret)
}