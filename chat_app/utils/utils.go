package utils

import (
	"errors"
	"fmt"
	
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPass(password string) (string){
	hash,err:=bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err!=nil{
		fmt.Println(err)
		return ""
	}
	return string(hash)
}

func CheckPasswordHash(password string,hash string){
	err:= bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	if err!=nil{
		fmt.Println("Invalid credentials",err)
		return
	}
	
}

type Claims struct{
	Id int `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateToken(userId int,user_name string) (string,error){
	claims:=Claims{
		Id: userId,
		Name: user_name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	Secret:=[]byte (os.Getenv("Secret"))
	return token.SignedString(Secret)
}

func ValidateToken(token string)(*Claims,error){
	claims:=&Claims{}

	parsedtoken,err:=jwt.ParseWithClaims(token,claims,func(t *jwt.Token) (interface{}, error){
		return []byte(os.Getenv("Secret")),nil
	})
	if err != nil || parsedtoken == nil || !parsedtoken.Valid {
			return nil,errors.New("Invalid Token")
		}

		return claims,err
}

