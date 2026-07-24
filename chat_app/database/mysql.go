package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"main.go/models"
	"main.go/utils"
)


var DB *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dsn := os.Getenv("MYSQL_DSN")

	// instead of checking connection for once
	// we will check continuosly
	var connErr error
	for i := 0; i < 10; i++ {
		DB, connErr = sql.Open("mysql", dsn) 
		if connErr == nil {
			connErr = DB.Ping()
			if connErr == nil {
				break
			}
		}
		log.Println("Waiting for MySQL...")
		time.Sleep(3 * time.Second)
	}
	if connErr != nil {
		log.Fatal("Could not connect to MySQL:", connErr)
	}
	fmt.Println("The connection is established")
}

func Existing(req models.Register) (bool,error){

	var id int
	err:=DB.QueryRow("select id from user where Email=?",req.Email).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // user does not exist
		}
		return false, err // real database error
	}

	return true,nil
}

func Register(user models.Register){
	result,err:=DB.Exec("insert into user (Name,Email,Password) values (?,?,?)",user.Name,user.Email,user.Password)
	if err!=nil{
		log.Fatal(err)
	}
	id,_:=result.LastInsertId()
	fmt.Println("The user id ->",id)
}

func Login(req models.Login) models.User{
var user models.User
	err:=DB.QueryRow("select * from user where email=?",req.Email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err!=nil{
		log.Fatal(err)
	}
	
	// compare Hashed password
	utils.CheckPasswordHash(req.Password,user.Password)
	return user
}

func SaveMessage(SenderId int ,ReceiverId int ,content string) (models.Message,error){
	result, err := DB.Exec(
		"insert into messages (sender_id, receiver_id, content) values (?, ?, ?)",
		SenderId, ReceiverId, content,
	)
	if err != nil {
		return models.Message{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.Message{}, err
	}
 
	var m models.Message
	err = DB.QueryRow(
		"select id, sender_id, receiver_id, content, created_at from messages where id = ?",
		id,
	).Scan(&m.ID, &m.SenderId, &m.ReceiverId, &m.Content, &m.CreatedAt)
	if err != nil {
		return models.Message{}, err
	}
	return m, nil
}