package main

import(
	"fmt"
	// mysql connector
	//"database/sql"
	"time"
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

func (lib *Library) insert() error{
	current_time := time.Now();
	last, _ := time.ParseDuration("-840h")
	time_4_1 := current_time.Add(last)
	last, _ = time.ParseDuration("-624h")
	time_4_10 := current_time.Add(last)
	last, _ = time.ParseDuration("-288h")
	time_4_24 := current_time.Add(last)
	last, _ = time.ParseDuration("-504h")
	time_4_15 := current_time.Add(last)
	last, _ = time.ParseDuration("-168h")
	time_4_29 := current_time.Add(last)
	last, _ = time.ParseDuration("-600h")
	time_4_11 := current_time.Add(last)
	last, _ = time.ParseDuration("-384h")
	time_4_20 := current_time.Add(last)
	last, _ = time.ParseDuration("-48h")
	time_5_4 := current_time.Add(last)
	last, _ = time.ParseDuration("-480h")
	time_4_16 := current_time.Add(last)
	last, _ = time.ParseDuration("-456h")
	time_4_17 := current_time.Add(last)
	last, _ = time.ParseDuration("-120h")
	time_5_1 := current_time.Add(last)
	last, _ = time.ParseDuration("-408h")
	time_4_19 := current_time.Add(last)
	last, _ = time.ParseDuration("-72h")
	time_5_3 := current_time.Add(last)
	last, _ = time.ParseDuration("-96h")
	time_5_2 := current_time.Add(last)
	last, _ = time.ParseDuration("216h")
	time_5_15 := current_time.Add(last)
	last, _ = time.ParseDuration("240h")
	time_5_16 := current_time.Add(last)
	last, _ = time.ParseDuration("264h")
	time_5_17 := current_time.Add(last)
	last, _ = time.ParseDuration("288h")
	time_5_18 := current_time.Add(last)
	
	
	lib.db.Exec("INSERT INTO STUDENT(user, password)VALUES('student1', '111'), ('student2', '222'),('student3', '333')")
	lib.db.Exec("INSERT INTO BOOK(title, author, ISBN, state)VALUES('book1', 'li', '100', 'borrowed'), ('book2', 'wang', '200', 'borrowed'), ('book3', 'zhang', '300', 'borrowed'), ('book4', 'sun', '400', 'borrowed')")
	_, err := lib.db.Exec(`INSERT INTO HISTORY(time, ISBN, user, op, due)
					VALUES(?, '200', 'student1', 'borrow', ?), 
						(?, '400', 'student2', 'borrow', ?), 
						(?, '100', 'student3', 'borrow', ?), 
						(?, '200', 'student1', 'return', ?),
						(?, '400', 'student2', 'extend', ?),
						(?, '100', 'student3', 'return', ?),
						(?, '100', 'student2', 'borrow', ?),
						(?, '200', 'student2', 'borrow', ?),
						(?, '300', 'student1', 'borrow', ?),
						(?, '300', 'student1', 'extend', ?),
						(?, '300', 'student1', 'extend', ?),
						(?, '300', 'student1', 'extend', ?)`, 
					time_4_1, time_4_15,
					time_4_10, time_4_24,
					time_4_15, time_4_29,
					time_4_11, time_4_11,
					time_4_20, time_5_4,
					time_4_16, time_4_16,
					time_4_17, time_5_1,
					time_4_19, time_5_3,
					time_5_1, time_5_15,
					time_5_2, time_5_16,
					time_5_3, time_5_17,
					time_5_4, time_5_18)
	if err != nil{
		panic(err)
	}	
	return nil
}

func main(){
	var lib *Library
	lib = new(Library)
	lib.ConnectDB()
	lib.insert()
}
