package books

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/books-list/driver"
	"main/books-list/models"
	bookRepository "main/books-list/repository/book"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var books []models.Book
var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// InitAPI initiates routes
func InitAPI() {
	db = driver.ConnectDB()
	router := mux.NewRouter()
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books", AddBook).Methods("POST")
	router.HandleFunc("/books", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", RemoveBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetBooks returns list of books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	books = []models.Book{}
	bookRepo := bookRepository.BookRepository{}
	books = bookRepo.GetBooks(db, book, books)

	json.NewEncoder(w).Encode(books)
}

// GetBook returns one book by id
func GetBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	params := mux.Vars(r)

	books = []models.Book{}
	bookRepo := bookRepository.BookRepository{}

	id, err := strconv.Atoi(params["id"])
	logFatal(err)
	book = bookRepo.GetBook(db, book, id)

	json.NewEncoder(w).Encode(book)
}

// AddBook adds new book to books list
func AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)
	bookRepo := bookRepository.BookRepository{}
	bookID = bookRepo.AddBook(db, book)

	json.NewEncoder(w).Encode(bookID)
}

// UpdateBook updates existing book by id
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	bookRepo := bookRepository.BookRepository{}
	rowsUpdated := bookRepo.UpdateBook(db, book)

	json.NewEncoder(w).Encode(rowsUpdated)
}

// RemoveBook deletes book by id
func RemoveBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	bookRepo := bookRepository.BookRepository{}
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	rowsDeleted := bookRepo.RemoveBook(db, id)

	json.NewEncoder(w).Encode(rowsDeleted)
}
