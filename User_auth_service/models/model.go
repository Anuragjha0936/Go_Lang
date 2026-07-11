package models

import "database/sql"

// import "database/sql"

// import "github.com/golang-jwt/jwt/v5"



type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Age sql.NullInt32 `json:"age"`
	Gender sql.NullString `json:"gender"`
	Leetcode_profile sql.NullString `json:"leetcode"`
	Password string `json:"password"`
}

type Cred struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

