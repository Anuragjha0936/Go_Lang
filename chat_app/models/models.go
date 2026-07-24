package models

import "time"

type User struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Register struct{
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Login struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type Message struct{
	ID int `json:"id"`
	SenderId int `json:"sender_id"`
	ReceiverId int `json:"receiver_id"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type WsMessage struct{
	ReceiverId int `json:"receiver_id"`
	Content string `json:"content"`
}