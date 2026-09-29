package controllers

import (
	"encoding/json"
	"net/http"
	"main.go/services"
)

type AIRequest struct{
	Message string `json:"message"`
}

type AIResponse struct{
	Reply string `json:"reply"`
}
func AiChat(w http.ResponseWriter,r *http.Request){
	var req AIRequest

	err:=json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w,"Invalid Request",http.StatusBadRequest)
		return
	}

	reply,err:=services.GetAiResponse(req.Message)

	if err!=nil{
		http.Error(w, "Failed to get AI response", http.StatusInternalServerError)
		return
	}

	response:=AIResponse{
		Reply:reply,
	}

	// telling the client that i am sending the response in the json format
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)
}