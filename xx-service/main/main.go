package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

var DB *sql.DB = setupGlobalDatabase()

func setupGlobalDatabase() *sql.DB {
	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}
	return db
}

type User struct {
	ID       int
	Name     string
	Language string
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		err2 := fmt.Errorf("Couldn't set up database: %w", err)
		return nil, err2
	}

	sqlStmt := `
	create table users (id integer not null primary key AUTOINCREMENT, name text, language text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		err2 := fmt.Errorf("Couldn't create table: %w", err)
		return nil, err2
	}

	_, err = db.Exec("insert into users(name, language) values('jason', 'ruby'), ('sue', 'python')")
	if err != nil {
		panic(err)
	}

	return db, nil
}

func main() {
	http.HandleFunc("/", Handler)
	fmt.Println("Listening on port http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("select * from users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Language)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	fmt.Fprintf(w, "users: %#v\n", users)
}
