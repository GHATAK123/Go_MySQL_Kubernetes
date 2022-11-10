package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const API_PATH = "/api/v1/book"

type library struct {
	dbHost, dbPass, dbName string
}

type Book struct {
	Id, Name, Author, Isbn string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "9570"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods("GET")
	r.HandleFunc(apiPath, l.createBooks).Methods("POST")
	http.ListenAndServe("localhost:8080", r)
}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("Error querying books %s\n", err.Error())
	}
	books := []Book{}
	for rows.Next() {
		var id, name, author, isbn string
		rows.Scan(&id, &name, &author, &isbn)
		if err != nil {
			log.Fatalf("While Scanning %s\n", err.Error())
		}
		aBook := Book{
			Id:     id,
			Name:   name,
			Author: author,
			Isbn:   isbn,
		}
		books = append(books, aBook)
	}
	json.NewEncoder(w).Encode(books)
	l.closeConnection(db)
}

func (l library) createBooks(w http.ResponseWriter, r *http.Request) {
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)
	db := l.openConnection()

	insertQuery, err := db.Prepare("insert into books value (?,?,?,?)")
	if err != nil {
		log.Fatalf("While Preparing the db Query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("While begining the db Query %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Author, book.Isbn)
	if err != nil {
		log.Fatalf("While executing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("While commit the db Query %s\n", err.Error())
	}

	l.closeConnection(db)

}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("While Opening the database %s\n", err.Error())
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("While Closing the database %s\n", err.Error())
	}
}
