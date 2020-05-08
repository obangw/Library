package main


import (
	"testing"
	"database/sql"
	"time"
)
const(
	TIME1 = "2020-04-30 13:33:04"
	TIME2 = "2020-05-19 13:33:04"
)

func TestConnectDB(t *testing.T) {
	lib := Library{}
	err := lib.ConnectDB()
	if err != nil {
		t.Errorf("can't connect")
	}
}

func TestAddbook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	var forcheck string
	tables := []struct {
		x, y, z string
	}{
		{"book5", "a5", "500"},
		{"book6", "a6", "600"},
		{"book7", "a7", "700"},
	}
	for _, table := range tables{
		err := lib.AddBook(table.x, table.y, table.z)
		if err != nil{
			t.Errorf("can't Addbook")
		}
		err = lib.db.QueryRow("SELECT title FROM BOOK WHERE ISBN = ?", table.z).Scan(&forcheck)
		if err == sql.ErrNoRows{
			t.Errorf("Addbook unsuccessfully")
		}		
		if err != nil{
			panic(err)
		}
		if forcheck != table.x{
			t.Errorf("Addbook was incorrect, got: %s, want: %s", forcheck, table.x)
		}
	}
}

func TestAddaccount(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	var forcheck string
	tables := []struct {
		x, y string
	}{
		{"s4", "444"},
		{"s5", "555"},
		{"s6", "666"},
	}
	for _, table := range tables{
		err := lib.addaccount(table.x, table.y)
		if err != nil{
			t.Errorf("can't addaccount")
		}
		err = lib.db.QueryRow("SELECT password FROM STUDENT WHERE user = ?", table.x).Scan(&forcheck)
		if err == sql.ErrNoRows{
			t.Errorf("addaccount unsuccessfully")
		}		
		if err != nil{
			panic(err)
		}
		if forcheck != table.y{
			t.Errorf("Addbook was incorrect, got: %s, want: %s", forcheck, table.y)
		}
	}
}

func TestSearch_title(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	tables := []struct {
		x book_info
	}{
		{book_info{"book1", "li", "100", "borrowed"}},
		{book_info{"book5", "a5", "500", "ready"}},
	}
	for _, table := range tables{
		title, author, ISBN, state , err := lib.search_title(table.x.title)
		if err != nil{
			t.Errorf("search unsuccessfully")
		}
		if title != table.x.title {
			t.Errorf("search_title was incorrect, got: %s, want: %s", title, table.x.title)
		}
		if author != table.x.author {
			t.Errorf("search_title was incorrect, got: %s, want: %s", author, table.x.author)
		}
		if ISBN != table.x.ISBN {
			t.Errorf("search_title was incorrect, got: %s, want: %s", ISBN, table.x.ISBN)
		}
		if state != table.x.state {
			t.Errorf("search_title was incorrect, got: %s, want: %s", state, table.x.state)
		}
	}
}

func TestSearch_author(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	tables := []struct {
		x book_info
	}{
		{book_info{"book6", "a6", "600", "ready"}},
		{book_info{"book5", "a5", "500", "ready"}},
	}
	for _, table := range tables{
		title, author, ISBN, state , err := lib.search_author(table.x.author)
		if err != nil{
			t.Errorf("search unsuccessfully")
		}
		if title != table.x.title {
			t.Errorf("search_title was incorrect, got: %s, want: %s", title, table.x.title)
		}
		if author != table.x.author {
			t.Errorf("search_title was incorrect, got: %s, want: %s", author, table.x.author)
		}
		if ISBN != table.x.ISBN {
			t.Errorf("search_title was incorrect, got: %s, want: %s", ISBN, table.x.ISBN)
		}
		if state != table.x.state {
			t.Errorf("search_title was incorrect, got: %s, want: %s", state, table.x.state)
		}
	}
}

func TestSearch_ISBN(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	tables := []struct {
		x book_info
	}{
		{book_info{"book1", "li", "100", "borrowed"}},
		{book_info{"book3", "zhang", "300", "borrowed"}},
	}
	for _, table := range tables{
		title, author, ISBN, state , err := lib.search_ISBN(table.x.ISBN)
		if err != nil{
			t.Errorf("search unsuccessfully")
		}
		if title != table.x.title {
			t.Errorf("search_title was incorrect, got: %s, want: %s", title, table.x.title)
		}
		if author != table.x.author {
			t.Errorf("search_title was incorrect, got: %s, want: %s", author, table.x.author)
		}
		if ISBN != table.x.ISBN {
			t.Errorf("search_title was incorrect, got: %s, want: %s", ISBN, table.x.ISBN)
		}
		if state != table.x.state {
			t.Errorf("search_title was incorrect, got: %s, want: %s", state, table.x.state)
		}
	}
}

func TestBorrow(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	var forcheck string
	nomean := time.Now().Format("2006-1-2 15:04:05")
	tables := []struct {
		x string
		y int
	}{
		{"100", 0},
		{"500", 1},
	}
	for _, table := range tables{
		result, _ := lib.borrow(table.x)
		if result != table.y{
			t.Errorf("Borrow was incorrect, got: %d, want: %d", result, table.y)
		}
		if result == 1{
			err := lib.db.QueryRow("SELECT state FROM BOOK WHERE ISBN = ?", table.x).Scan(&forcheck)
			if err != nil{
				panic(err)
			}
			if forcheck != "borrowed"{
				t.Errorf("Error! Book state didn't changed")
			}
			err = lib.db.QueryRow("SELECT MAX(time), op FROM HISTORY WHERE ISBN = ? GROUP BY time, user, ISBN ORDER BY time DESC", table.x).Scan(&nomean, &forcheck)
			if err != nil{
				panic(err)
			}
			if forcheck != "borrow"{
				t.Errorf("Error! History didn't changed")
			}
		}
	}
}

func TestCheck_account(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	tables := []struct {
		x string
		y int
	}{
		{"student2", 3},
		{"student1", 0},
	}
	for _, table := range tables{
		result := lib.check_account(table.x, 0)
		if result != table.y{
			t.Errorf("Error! Checkoverdue error!")
		}
	}
}

func TestRemove(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	var forcheck string
	nomean := time.Now().Format("2006-1-2 15:04:05")
	tables := []struct {
		x string
		y string
	}{
		{"600", "forget"},
	}
	for _, table := range tables{
		lib.Remove(table.x, table.y)
		err := lib.db.QueryRow("SELECT title FROM BOOK WHERE ISBN = ?", table.x).Scan(&forcheck)
		if err != sql.ErrNoRows{
			t.Errorf("The %s didn't delete", table.x)
		}
		err = lib.db.QueryRow("SELECT MAX(time), op FROM HISTORY WHERE ISBN = ? GROUP BY time, user, ISBN ORDER BY time DESC", table.x).Scan(&nomean, &forcheck)
		if err != nil{
			panic(err)
		}
		if forcheck != "other"{
			t.Errorf("ERROR! Remove didn't insert into table history!")
		}
	}
}

func TestHistory(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	tables := []struct{
		x string
		y error
	}{
		{"student3", nil},
		{"s5", sql.ErrNoRows},
	}
	for _, table := range tables{
		err := lib.history(table.x)
		if err != table.y{
			t.Errorf("ERROR! History didn't return successfully!")
		}
	}
}

func TestRemain(t *testing.T){
	tables := []struct{
		x string
	}{
		{"student3"},
		{"student1"},
	}
	lib := Library{}
	lib.ConnectDB()
	for _, table := range tables{
		forcheck, _ := lib.remain(table.x)
		if table.x == "student3"{
			if forcheck != nil{
				t.Errorf("ERROR! %s isn't correct", table.x)
			}
		}else{
			if forcheck[0] != "300"{
				t.Errorf("ERROR! %s isn't correct", table.x)
			}
		}
	}
}

func TestCheckddl(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	tables := []struct{
		x, y, z string
	}{
		{"student3", "100", TIME1},
		{"student1", "300", TIME2},
	}
	for _, table := range tables{
		result, _ := lib.checkddl(table.x, table.y)
		if result != table.z{
			t.Errorf("ERROR! checkddl %s, %s, got: %s, expect: %s", table.x, table.y, result, table.z)
		}
	}
}

func TestExtend(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.user = "student2"
	var usercheck string
	nomean := time.Now().Format("2006-1-2 15:04:05")
	tables := []struct{
		x string
	}{
		{"200"},
		{"400"},
	}
	for _, table := range tables{
		err := lib.extend(table.x)
		err = lib.db.QueryRow("SELECT MAX(time), user FROM HISTORY WHERE ISBN = ? AND op = 'extend' GROUP BY ISBN, time, user ORDER BY time DESC", table.x).Scan(&nomean, &usercheck)
		if err != nil{
			panic(err)
		}
		if usercheck != "student2"{
			t.Errorf("ERROR! Extend unsuccessfully!")
		}
	}
}

func TestReturnbook(t *testing.T){
	lib := Library{}
	lib.ConnectDB()
	lib.user = "student2"
	var usercheck string
	nomean := time.Now().Format("2006-1-2 15:04:05")
	tables := []struct{
		x string
	}{
		{"100"},
	}
	for _, table := range tables{
		err := lib.returnbook(table.x)
		err = lib.db.QueryRow("SELECT MAX(time), user FROM HISTORY WHERE ISBN = ? AND op = 'return' GROUP BY ISBN, time, user ORDER BY time DESC", table.x).Scan(&nomean, &usercheck)
		if err != nil{
			panic(err)
		}
		if usercheck != "student2"{
			t.Errorf("ERROR! return unsuccessfully!")
		}
		err = lib.db.QueryRow("SELECT state FROM BOOK WHERE ISBN = ?", table.x).Scan(&usercheck)
		if err != nil{
			panic(err)
		}
		if usercheck != "ready"{
			t.Errorf("ERROR! return unsuccessfully!")
		}
	}
}
