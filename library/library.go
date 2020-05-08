package main

import (
	"fmt"
	"time"
	// mysql connector
	"database/sql"
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

type book_info struct{
	title, author, ISBN, state string
}

type history_info struct{
	time, ISBN, user, op, due string
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

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) error {
	_, err := lib.db.Exec("INSERT INTO BOOK(title, author, ISBN, state)VALUES(?,?,?,?)", title, author, ISBN, "ready")
	if err != nil{
		panic(err)
	}
	return err
}

func (lib *Library) welcome() int{
	var password_saved, user, password string
	fmt.Println("Welcome to the Library Management System!")
	fmt.Println("Please login first.End input with [ENTER]")
	fmt.Printf("Username:  ")
	fmt.Scanf("%s", &user)
	row, err := lib.db.Query("SELECT Password FROM STUDENT WHERE user = ?", user)
	if err != nil{
		panic(err)
	}	
	for row.Next(){	
		err = row.Scan(&password_saved);
	}
	if err != nil{
		panic(err)
	}	
	row.Close()
	for i := 3 ; i > 0; i--{
		fmt.Printf("Password:  ")
		fmt.Scanf("%s", &password)
		if password != password_saved{
			fmt.Printf("Wrong Password. You have %d times left\n", i-1)
		}else{
			break
		}
	}
	if password != password_saved{
		return 1
	}
	lib.user = user
	return 0
}

func (lib *Library) addaccount(user, password string) error{
	_, err := lib.db.Exec("INSERT INTO STUDENT(user, password)VALUE(?,?)", user, password)
	if err != nil{panic(err)}
	fmt.Println("DONE")
	return nil
}

func (lib *Library) search_title(title string) (string, string, string, string, error){
	var for_check string
	var result book_info
	check := lib.db.QueryRow("SELECT title FROM BOOK WHERE title = ?", title).Scan(&for_check)
	if check == sql.ErrNoRows{
		fmt.Println("No information")
		fmt.Println("DONE")
		return result.title, result.author, result.ISBN, result.state, nil
	}else if check != nil{panic(check)}
	row, err := lib.db.Query("SELECT title, author, ISBN, state FROM BOOK WHERE title = ?", title)
	if err != nil{panic(err)}
	for row.Next(){
		err = row.Scan(&result.title, &result.author, &result.ISBN, &result.state)
		if err != nil{
			panic(err)
		}	
		fmt.Println(result)
	}
	row.Close()
	fmt.Println("DONE")
	return result.title, result.author, result.ISBN, result.state, nil
}

func (lib *Library) search_author(author string) (string, string, string ,string, error){
	var for_check string
	var result book_info
	check := lib.db.QueryRow("SELECT title FROM BOOK WHERE author = ?", author).Scan(&for_check)
	if check == sql.ErrNoRows{
		fmt.Println("No information")
		fmt.Println("DONE")
		return result.title, result.author, result.ISBN, result.state, nil
	}else if check != nil{panic(check)}
	row, err := lib.db.Query("SELECT title, author, ISBN, state FROM BOOK WHERE author = ?", author)
	if err != nil{panic(err)}
	for row.Next(){
		err = row.Scan(&result.title, &result.author, &result.ISBN, &result.state)
		if err != nil{
			panic(err)
		}	
		fmt.Println(result)
	}
	row.Close()
	fmt.Println("DONE")
	return result.title, result.author, result.ISBN, result.state, nil
}

func (lib *Library) search_ISBN(ISBN string) (string, string, string ,string, error){
	var for_check string
	var result book_info
	check := lib.db.QueryRow("SELECT title FROM BOOK WHERE ISBN = ?", ISBN).Scan(&for_check)
	if check == sql.ErrNoRows{
		fmt.Println("No information")
		fmt.Println("DONE")
		return result.title, result.author, result.ISBN, result.state, nil
	}else if check != nil{panic(check)}
	row, err := lib.db.Query("SELECT title, author, ISBN, state FROM BOOK WHERE ISBN = ?", ISBN)
	if err != nil{panic(err)}
	for row.Next(){
		err = row.Scan(&result.title, &result.author, &result.ISBN, &result.state)
		if err != nil{
			panic(err)
		}	
		fmt.Println(result)
	}
	row.Close()
	fmt.Println("DONE")
	return result.title, result.author, result.ISBN, result.state, nil
}

func (lib *Library) borrow(ISBN string) (int, error){
	var state_temp string	
	err := lib.db.QueryRow("SELECT state FROM BOOK WHERE ISBN = ?", ISBN).Scan(&state_temp)
	if err != nil{panic(err)}
	if state_temp != "ready"{
		fmt.Println("Error! This book has been borrowed by others")		
		return 0, nil
	}
	current_time := time.Now()
	last , _ := time.ParseDuration("336h")
	duetime := current_time.Add(last)
	if err != nil{return 0, err}
	if lib.overdue_num >= 3{
		fmt.Println("You can't borrow a book until you return a book first")
		return 0, nil
	}
	_, err = lib.db.Exec("UPDATE BOOK SET state = ? WHERE ISBN = ?", "borrowed", ISBN)
	if err != nil{panic(err)}
	_, err = lib.db.Exec("INSERT INTO HISTORY(time, ISBN, user, op, due) VALUE(?,?,?,?,?)", current_time.Format("2006-1-2 15:04:05"), ISBN, lib.user, "borrow", duetime.Format("2006-1-2 15:04:05"))
	if err != nil{panic(err)}
	fmt.Println("DONE")
	return 1, nil
}

func (lib *Library) check_account(user string, tag int) int{
	count := 0
	var result book_info
	current_time := time.Now()
	row, err := lib.db.Query(`WITH 	BORROW(time, duetime, ISBN, user) AS 
						(
						SELECT MAX(time), MAX(due), ISBN, user 
						FROM HISTORY 
						WHERE user = ? AND (op = 'borrow' OR op = 'extend') 
						GROUP BY ISBN, user, time
						)
				SELECT DISTINCT BORROW.ISBN 
				FROM BORROW, HISTORY
				WHERE ((BORROW.ISBN = HISTORY.ISBN AND (HISTORY.op = 'return' OR HISTORY.op = 'other') AND HISTORY.user = ? AND BORROW.time > HISTORY.time) OR (BORROW.ISBN NOT IN (SELECT ISBN FROM HISTORY WHERE (HISTORY.op = 'return' OR HISTORY.op = 'other') AND HISTORY.user = ?))) AND BORROW.duetime < ?`, user, user, user, current_time.Format("2006-1-2 15:04:05"))
	if err != nil{
			panic(err)
		}
	for row.Next(){
		row.Scan(&result.ISBN)
		count ++
		err = lib.db.QueryRow("SELECT title, author, state FROM BOOK WHERE ISBN = ?", result.ISBN).Scan(&result.title, &result.author, &result.state)
		if err != nil{
			panic(err)
		}
		if tag == 1 {
			fmt.Println(result)
		}
	}
	row.Close()
	return count
}
// etc...
func command_fault(){
	fmt.Println("[ERROR]command not found")
	fmt.Println("Input help for more information")
}

func (lib *Library) Remove(ISBN, explanation string) error {
	current_time := time.Now()
	_, err := lib.db.Exec("INSERT INTO HISTORY(time, ISBN, user, op, due) VALUE(?,?,?,?,?)", current_time.Format("2006-1-2 15:04:05"), ISBN, lib.user, "other", current_time.Format("2006-1-2 15:04:05"))
	if err != nil{panic(err)}
	_, err = lib.db.Exec("DELETE FROM BOOK WHERE ISBN = ?", ISBN)
	if err != nil{panic(err)}
	return nil;
}

func checkerror(err error){
	if err == nil{
		fmt.Println("DONE")
	}else{
		fmt.Println(err)
	}
}

func (lib *Library) history(user string) error{
	var for_check string	
	check := lib.db.QueryRow("SELECT op FROM HISTORY WHERE user = ? ORDER BY time", user).Scan(&for_check)
	if check == sql.ErrNoRows{
		fmt.Println("No information")
		return check
	}else if check != nil{panic(check)}
	row, err := lib.db.Query("SELECT time, ISBN, user, op, due FROM HISTORY WHERE user = ? ORDER BY time", user)
	if err != nil{
		panic(check)
	}
	var result history_info	
	for row.Next(){
		row.Scan(&result.time, &result.ISBN, &result.user, &result.op, &result.due)
		fmt.Println(result)
	}
	row.Close()
	return nil
}
func (lib *Library) remain(user string) ([]string, error){
	var for_check string
	var check error
	var output []string
	check = lib.db.QueryRow(`WITH BORROW(time, ISBN) AS (
						SELECT MAX(time), ISBN 
						FROM HISTORY 
						WHERE user = ? AND op = 'borrow' 
						GROUP BY ISBN, user, time
						)
				SELECT BOOK.title
				FROM BORROW, HISTORY, BOOK 
				WHERE BOOK.ISBN = BORROW.ISBN AND 
				((BORROW.ISBN NOT IN (SELECT ISBN FROM HISTORY WHERE (HISTORY.op = 'return' OR HISTORY.op = 'other') AND HISTORY.user = ?)) OR (BORROW.ISBN = HISTORY.ISBN AND HISTORY.user = ? AND (HISTORY.op = 'return' OR HISTORY.op = 'other') AND BORROW.time > HISTORY.time))`, user, user, user).Scan(&for_check)
	if check == sql.ErrNoRows{
		fmt.Println("No information")
		return output, nil
	}else if check != nil{panic(check)}
	var result book_info
	row, err := lib.db.Query(`WITH BORROW(time, ISBN) AS (
						SELECT MAX(time), ISBN 
						FROM HISTORY 
						WHERE user = ? AND op = 'borrow' 
						GROUP BY ISBN, user, time
						)
				SELECT DISTINCT BOOK.title, BOOK.ISBN, BOOK.author, BOOK.state
				FROM BORROW, HISTORY, BOOK 
				WHERE BOOK.ISBN = BORROW.ISBN AND 
				((BORROW.ISBN NOT IN (SELECT ISBN FROM HISTORY WHERE (HISTORY.op = 'return' OR HISTORY.op = 'other') AND HISTORY.user = ?)) OR (BORROW.ISBN = HISTORY.ISBN AND HISTORY.user = ? AND (HISTORY.op = 'return' OR HISTORY.op = 'other') AND BORROW.time > HISTORY.time))`, user, user, user)
	if err != nil{
		panic(err)
	}	
	for row.Next(){
		row.Scan(&result.title, &result.ISBN, &result.author, &result.state)
		output = append(output, result.ISBN)
		fmt.Println(result)
	}
	row.Close()
	return output, nil
}

func (lib *Library) checkddl(user, ISBN string) (string, error){
	var for_check string
	result := time.Now().Format("2006-1-2 15:04:05")
	check := lib.db.QueryRow("SELECT ISBN FROM HISTORY WHERE user = ? AND (op = ? OR op = ?) AND ISBN = ?", user, "borrow", "extend", ISBN).Scan(&for_check)
	if check == sql.ErrNoRows{
		fmt.Println("No information")
		return result, nil
	}else if check != nil{panic(check)}
	err:= lib.db.QueryRow("SELECT MAX(due) FROM HISTORY WHERE user = ? AND (op = ? OR op = ?) AND ISBN = ?", user, "borrow", "extend", ISBN).Scan(&result)
	if err != nil{
		return result, err
	}
	fmt.Println(result)
	return result, nil
}

func (lib *Library) extend(ISBN string) error{
	var times int
	err := lib.db.QueryRow("SELECT COUNT(time) FROM HISTORY WHERE user = ? AND ISBN = ? AND op = ? AND time > (SELECT MAX(time) FROM HISTORY WHERE user = ? AND ISBN = ? AND op = ?)", lib.user, ISBN, "extend", lib.user, ISBN, "BORROW").Scan(&times)	
	if (err != nil) && (err != sql.ErrNoRows) {
		panic(err)
	}
	if times == 3 {
		fmt.Println("You have already extended 3 times!")
		return nil
	}else{
		state := lib.make_sure(lib.user, ISBN)
		if state == 1{
			fmt.Println("You havn't borrowed this book yet!")
			return nil
		}	
		fmt.Printf("You have %d times left\n", 2-times)
		//Make sure this people has borrowed the book
		current_time := time.Now()
		last , _ := time.ParseDuration("336h")
		duetime := current_time.Add(last)
		_, err := lib.db.Exec("INSERT INTO HISTORY(time, ISBN, user, op, due) VALUE(?,?,?,?,?)", current_time.Format("2006-1-2 15:04:05"), ISBN, lib.user, "extend", duetime.Format("2006-1-2 15:04:05"))
		if err != nil{
			panic(err)
		}
	}
	return nil
}
func (lib *Library)make_sure(user, ISBN string) int{
	for_check := time.Now().Format("2006-1-2 15:04:05")
	var user_output, state string
	err := lib.db.QueryRow("SELECT MAX(time), user FROM HISTORY WHERE ISBN = ? AND op = 'borrow' GROUP BY user, time, ISBN ORDER BY time DESC", ISBN).Scan(&for_check, &user_output)
	if err == sql.ErrNoRows{
		return 1
	}else if err != nil{
			panic(err)
		}
	if user != user_output{
		return 1
	}
	err = lib.db.QueryRow("SELECT state FROM BOOK WHERE ISBN = ?", ISBN).Scan(&state)
	if err == sql.ErrNoRows{
		return 1
	}else if err != nil{
			panic(err)
		}
	if state != "borrowed" {
		return 1
	}
	return 0
}

func (lib *Library) returnbook(ISBN string) error{
	state_temp := lib.make_sure(lib.user, ISBN)
	if state_temp == 1{
		fmt.Println("You havn't borrow this book yet")
		return nil
	}else{
		_, err := lib.db.Exec("UPDATE BOOK SET state = ? WHERE ISBN = ?", "ready", ISBN)
		if err != nil{
			panic(err)
		}
		current_time := time.Now()
		_, err = lib.db.Exec("INSERT INTO HISTORY(time, ISBN, user, op, due) VALUE(?,?,?,?,?)", current_time.Format("2006-1-2 15:04:05"), ISBN, lib.user, "return", current_time.Format("2006-1-2 15:04:05"))
		if err != nil{
			return err
		}
	}
	fmt.Println("DONE")
	return nil
}
func help(){
	fmt.Println("The command with [] should be filled with correct information")	
	fmt.Println("addbook [title] [author] [ISBN]			FOR ADMINSTRATOR ONLY. Add book to the library")
	fmt.Println("remove [title] [author] [ISBN] [explanation]	If you can't return a book, you can use this function")
	fmt.Println("addaccount	[newuser] [newpassword]			FOR ADMINISTRATOR ONLY. Add account")
	fmt.Println("search						Search a book")
	fmt.Println("borrow [ISBN]					Borrow a book with ISBN")
	fmt.Println("history 						Show the operator history of a account.")
	fmt.Println("remain						Show the book remained to be returned")
	fmt.Println("checkddl [ISBN]					Check the book's returning ddl with ISBN")
	fmt.Println("extend [ISBN]					Extend the book's returning ddl. If you have extended three times, you can't use this command")
	fmt.Println("checkoverdue					Check if you have any overdue book")
	fmt.Println("return [ISBN]					Return a book with ISBN")
	fmt.Println("quit						Quit")
	fmt.Println(" ")
	fmt.Println(" ")
}

func main() {
	var user, command, title, author, ISBN, explanation, newuser, newpassword, newpassword_again string
	var lib *Library
	lib = new(Library)
	var way int
	running := 1
	lib.ConnectDB()
	judge := lib.welcome()
	if judge == 1{return}
	fmt.Printf("Welcome, %s\n", lib.user)
	lib.overdue_num = lib.check_account(lib.user, 0)
	for running == 1 {
		lib.overdue_num = lib.check_account(lib.user, 0)
		fmt.Scanf("%s", &command)
		switch command{
			case "addbook" :{	if lib.user != User{
							fmt.Println("No Authority")
						}else{
						fmt.Scanf("%s%s%s", &title, &author, &ISBN)
						err := lib.AddBook(title, author, ISBN)
						checkerror(err)
						}
	        	}
			case "remove" :	{	fmt.Scanf("%s%s%s%s", &title, &author, &ISBN, &explanation)
				 		err := lib.Remove(ISBN, explanation)
				 		checkerror(err)
			}
			case "addaccount":{	if lib.user != User {
							fmt.Println("No Authority")
						}else{
				   			fmt.Scanf("%s%s", &newuser, &newpassword)
							fmt.Printf("Confirm your password: ")
							fmt.Scanf("%s", &newpassword_again)
							if newpassword_again == newpassword {
								lib.addaccount(newuser, newpassword)
							}else{
								fmt.Println("Different password")
							 }	
				   		}
			}
			case "search":{		fmt.Println("BY 1.title  2.author  3.ISBN")
			  			fmt.Scanf("%d", &way)
						switch way{
							case 1: {fmt.Println("ATTENTION! Space is not allowed")
								 fmt.Println(" ")
								 fmt.Printf("Input the title : ")
								 fmt.Scanf("%s", &title)
								 lib.search_title(title)
								}
							case 2: {fmt.Println("ATTENTION! Space is not allowed")
								 fmt.Println(" ")
								 fmt.Printf("Input the author : ")
								 fmt.Scanf("%s", &author)
								 lib.search_author(author)
								}
							case 3: {fmt.Println("ATTENTION! Space is not allowed")
								 fmt.Println(" ")
								 fmt.Printf("Input the ISBN : ")
								 fmt.Scanf("%s", &ISBN)
								 lib.search_ISBN(ISBN)
								}
							default: {fmt.Println("Invalid option")}
						}		
			}
			case "borrow":	{ 	fmt.Scanf("%s", &ISBN)
						lib.borrow(ISBN)
		        }
			case "history":	{	if lib.user != User{
							lib.history(lib.user)
						}else{
							fmt.Printf("Input the user: ")
							fmt.Scanf("%s", &user)
							lib.history(user)
						}
			}
			case "remain":	{	if lib.user != User{
							lib.remain(lib.user)
						}else{
							fmt.Printf("Input the user: ")
							fmt.Scanf("%s", &user)
							lib.remain(user)
						}
			}
			case "checkddl":{	if lib.user != User{
							fmt.Printf("Input the book ISBN: ")
							fmt.Scanf("%s", &ISBN)
							lib.checkddl(lib.user, ISBN)
						}else{
							fmt.Println("Input the user: ")
							fmt.Scanf("%s", &user)
							fmt.Println("Input the ISBN: ")
							fmt.Scanf("%s", &ISBN)
							lib.checkddl(user, ISBN)
						}
			}
			case "extend":{		fmt.Scanf("%s", &ISBN)
						lib.extend(ISBN)
			}
			case "checkoverdue":{	lib.check_account(lib.user,1)
			}
			case "return":{		fmt.Scanf("%s", &ISBN)
						lib.returnbook(ISBN)
			}
			case "help":{
						help()
			}
			case "quit": 		running = 0
			default: 		command_fault()
		}
	}
}
