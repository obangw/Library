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

func (lib *Library) drop() error{
	lib.db.Exec("DROP TABLE STUDENT CASCADE")
	lib.db.Exec("DROP TABLE HISTORY CASCADE")
	lib.db.Exec("DROP TABLE BOOK CASCADE")
	return nil
}
func main(){
	var lib *Library
	lib = new(Library)
	lib.ConnectDB()
	lib.drop()
}
