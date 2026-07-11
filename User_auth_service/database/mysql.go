package database

import (
	"database/sql"
	"os"
	"time"

	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"main.go/models"
)
var DB *sql.DB
func init() {
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}
	dsn:=os.Getenv("MYSQL_DSN")
	// instead of checking connection for once 
	// we will check continuosly 
	for i := 0; i < 10; i++ {
    db, err := sql.Open("mysql", dsn)
    if err == nil {
        err = db.Ping()
        if err == nil {
            break
        }
    }

    log.Println("Waiting for MySQL...")
    time.Sleep(3 * time.Second)
}

if err != nil {
    log.Fatal("Could not connect to MySQL:", err)
}
fmt.Println("The connection is established")
}

func Login(req models.Cred) models.User{

	var user models.User
	err:=DB.QueryRow("select id,name,email,password from user where email=?",req.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err!=nil{
		log.Fatal(err)
	}

	// compare hashed password
	err=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password))
	if err!=nil{
		log.Fatal(err)
	}
	return user
}

func Register(req models.User){
	result,err:=DB.Exec("insert into user (name,email,password) values(?,?,?)",req.Name,req.Email,req.Password)
	if err!=nil{
		log.Fatal(err)
	}
	id,_:=result.LastInsertId()
	// when we insert a note in the db it return a unique id
	fmt.Println("the unique id ",id)
}

func Complete_p(id int,req models.User){
	_, err := DB.Exec(
		"update user set Age=?, Gender=?, Leetcode_profile=? where id=?",
		req.Age, req.Gender, req.Leetcode_profile, id,
	)
	if err!=nil{
		log.Fatal(err)
	}
	
}

func View_p(id int) models.User{
	var user models.User
	query:="select id,name,email,age,gender,Leetcode_profile from user where id=?"
	err:=DB.QueryRow(query,id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Gender,
		&user.Leetcode_profile,
	)
	if err!=nil{
		log.Fatal(err)
	}
	return user
}
