package controllers

import (
	"encoding/json"
	
	"net/http"
	"main.go/database"
	"main.go/models"
	"main.go/utils"
)

// Register
func Register(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var req models.Register

	err:=json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
	}

	if req.Email=="" || req.Name=="" || req.Password==""{
		http.Error(w, "Username, email and password are required", http.StatusBadRequest)
        return
	}

	// checking the user alreday exist or not with same mail
	exist,err:=database.Existing(req)
	if err!=nil{
		http.Error(w,"Internal Server Error",http.StatusInternalServerError)
		return
	}
	if exist{
		http.Error(w, "User already exists", http.StatusConflict)
	    return
	}

	// hash the password and then store it
	hashedPass:=utils.GenerateHashPass(req.Password)
	req.Password=hashedPass
	database.Register(req)
	
	json.NewEncoder(w).Encode(map[string]string{
	"message": "User successfully registered",
})
}

// Login
func Login(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")

	var req models.Login
	err:=json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email=="" || req.Password==""{
		http.Error(w,"Invalid username or Password",http.StatusBadRequest)
		return
	}

	user:=database.Login(req)
	
	// Generate Token
	token,err:=utils.GenerateToken(user.Id,user.Name)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token":token,
	})

}