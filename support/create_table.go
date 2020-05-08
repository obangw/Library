package main

import(
	"fmt"
	// mysql connector
	//"database/sql"
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "123"
	DBName   = "ass3"
)

type Library struct {
	user string
	overdue_num int
	db *sqlx.DB
}


func (lib *Library) ConnectDB() error {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	/*row, err := db.Query("SELECT user FROM STUDENT")
	if row.Next(){
	}else{
		db.Exec("INSERT INTO STUDENT(user, password) VALUE(?,?)", User, Password)
	}*/
	lib.db = db
	return nil
}

func (lib *Library) CreateTables() error {
	_, err := lib.db.Exec("CREATE TABLE STUDENT (user VARCHAR(40) NOT NULL, password VARCHAR(40) NOT NULL, PRIMARY KEY(user))")
	_, err = lib.db.Exec("CREATE TABLE HISTORY (time DATETIME NOT NULL, ISBN VARCHAR(13) NOT NULL, user VARCHAR(40) NOT NULL, op VARCHAR(6) NOT NULL, due DATETIME NOT NULL, PRIMARY KEY(time, user, ISBN))")
	if err != nil{
		panic(err)
	}	
	_, err = lib.db.Exec("CREATE TABLE BOOK (title VARCHAR(40) NOT NULL, author VARCHAR(40) NOT NULL, ISBN VARCHAR(40) NOT NULL, state VARCHAR(8) NOT NULL, PRIMARY KEY(ISBN))")
	if err != nil{
		panic(err)
	}
	lib.db.Exec("INSERT INTO STUDENT(user, password) VALUE(?,?)", User, Password)	
	return nil
}

func main(){
	var lib *Library
	lib = new(Library)
	lib.ConnectDB()
	lib.CreateTables()
}
