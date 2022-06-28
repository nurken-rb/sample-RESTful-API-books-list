package handlers

import (
	"books-list/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct{}

var books []models.Book

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//--------------------- Get list of all books ---------------------------

func (h Handler) Getbooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		books = []models.Book{}

		rows, err := db.Query("select * from books")
		logFatal(err)
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
			logFatal(err)
			books = append(books, book)

			json.NewEncoder(w).Encode(books)
		}
	}
}

//--------------------- Get book by id -----------------------------------

func (h Handler) Getbook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		params := mux.Vars(r)

		rows := db.QueryRow("select * from books where id=$1", params["id"])

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		json.NewEncoder(w).Encode(book)
	}
}

//---------------------- Add book ------------------------------------------

func (h Handler) Addbook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		var bookID int

		json.NewDecoder(r.Body).Decode(&book)

		err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;",
			book.Title, book.Author, book.Year).Scan(&bookID)
		logFatal(err)

		json.NewEncoder(w).Encode(bookID)
	}
}

// --------------------- Update book information ----------------------------

func (h Handler) Updatebook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)

		res, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
			&book.Title, &book.Author, &book.Year, &book.ID)
		logFatal(err)

		rowsUpdated, err := res.RowsAffected()
		logFatal(err)

		json.NewEncoder(w).Encode(rowsUpdated)
	}

}

//---------------------- Remove book -------------------------------------------

func (h Handler) Removebook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		res, err := db.Exec("delete from books where id = &1", params["id"])
		logFatal(err)

		rowsDeleted, err := res.RowsAffected()
		logFatal(err)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
