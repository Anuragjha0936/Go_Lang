package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"main.go/database"
	"main.go/middleware"
	"main.go/models"
)

func GetUser(w http.ResponseWriter,r *http.Request){

	currentUserId:=r.Context().Value(middleware.UserIDKey).(int)

	// when we have to return multiple rows we use Query
	rows,err:=database.DB.Query("select id,Name,Email from user where id!=?",currentUserId)
		if err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
         return
        }
	
		// always close the rows otherwise we will leak the underlying connection
		defer rows.Close()

		response:=make([]map[string]interface{},0)

		for rows.Next(){
			var user models.User

			err:=rows.Scan(
				&user.Id,
				&user.Name,
				&user.Email,
			)
			if err!=nil{
				http.Error(w,"Unable to read the user row",http.StatusInternalServerError)
				return
			}
			response=append(response,map[string]interface{}{
				"id":user.Id,
				"name":user.Name,
				"email":user.Email,
			})
		}
			// Check for errors encountered during iteration itself
			
			if rows.Err()!=nil{
				http.Error(w,"Error rading user",http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(response)
}

func GetMessage(w http.ResponseWriter,r *http.Request){
	currentUserId:=r.Context().Value(middleware.UserIDKey).(int)

	vars:=mux.Vars(r)
	otherId,err:=strconv.Atoi(vars["userId"])
	if err!=nil{
		http.Error(w,"Invalid user id",http.StatusInternalServerError)
		return
	}

	rows,err:=database.DB.Query(
		`select id,sender_id,receiver_id,content,created_at from messages WHERE (sender_id = ? AND receiver_id = ?) 
		OR (sender_id = ? AND receiver_id = ?) ORDER BY created_at ASC`,currentUserId,otherId,otherId,currentUserId,
	)

	if err!=nil{
		http.Error(w,"Failed to fetch message",http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages:=make([]models.Message,0)

	for rows.Next(){
		var m models.Message

		err:=rows.Scan(&m.ID,&m.SenderId,&m.ReceiverId,&m.Content,&m.CreatedAt)
		if err!=nil{
			http.Error(w,"unable to read message row",http.StatusInternalServerError)
			return
		}
		messages=append(messages, m)	
	}

	if rows.Err()!=nil{
		http.Error(w,"Error Reading Messages",http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

