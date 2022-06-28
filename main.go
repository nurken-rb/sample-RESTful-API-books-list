package main

import (
	"database/sql"
	"log"
	"net/http"

	"books-list/driver"
	"books-list/handlers"
	"books-list/models"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

//--------------------- Main function --------------------------------

func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()

	handler := handlers.Handler{}

	router.HandleFunc("/books", handler.Getbooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", handler.Getbook(db)).Methods("GET")
	router.HandleFunc("/books", handler.Addbook(db)).Methods("POST")
	router.HandleFunc("/books", handler.Updatebook(db)).Methods("PUT")
	router.HandleFunc("/books", handler.Removebook(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
