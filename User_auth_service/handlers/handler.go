package handlers

import (
	"encoding/json"
	"net/http"
	
	
	"golang.org/x/crypto/bcrypt"
	"main.go/database"
	
	"main.go/models"
	"main.go/utils"
)

// login handler
func Login(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Method","POST")

	// get the email and password from request body
	var req models.Cred

	err:=json.NewDecoder(r.Body).Decode(&req)

	
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]string{
				"error":"Invalid request body",
			},
		)
	}

	if req.Email == "" || req.Password == "" {
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "Email and Password are required",
    })
    return
}
	var user models.User
	
	
	user=database.Login(req)

	token,err:=utils.Generate_token(user.ID)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Token not generated")
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"token":token,
		},
	)
	
}

func Register(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Method","POST")

	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	hash,_:=bcrypt.GenerateFromPassword([]byte(user.Password),12)

	user.Password=string(hash)

	database.Register(user)
	json.NewEncoder(w).Encode("User is successfully Registered")
}

func View_profile(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Method","GET")
	
	userId:=r.Context().Value("id").(int)
	var user models.User
	user=database.View_p(userId)
	json.NewEncoder(w).Encode(user)
}
func Complete_Profile(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Method", "POST")
 
	userId, ok := r.Context().Value("id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}
 
	var data models.User
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	database.Complete_p(userId,data)
	json.NewEncoder(w).Encode("Profile Completed!!")

}
