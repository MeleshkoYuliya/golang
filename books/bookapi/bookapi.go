package bookapi

import (
	"database/sql"
	"encoding/json"
	"log"

	"net/http"
	"strconv"

	"github.com/MeleshkoYuliya/golang/books/repository"
	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/MeleshkoYuliya/golang/common/models"
	"github.com/MeleshkoYuliya/golang/notifier/notifierapi"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type booksService struct {
	db          *sql.DB
	books       []models.Book
	subscribers []models.Subscriber
}

var s booksService

func logFatal(err error) {
	if err != nil {
		spew.Dump(err)
	}
}

// InitAPI initiates routes
func InitAPI() {
	s.db = driver.GetDB()
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
	ctx := r.Context()

	bookRepo := repository.BookRepository{}
	books, err := bookRepo.GetBooks(ctx)
	logFatal(err)
	json.NewEncoder(w).Encode(books)
}

// GetBook returns one book by id
func GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	bookRepo := repository.BookRepository{}
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	book, err := bookRepo.GetBook(ctx, id)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}

// AddBook adds new book to books list
func AddBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)
	bookRepo := repository.BookRepository{}
	bookID, err := bookRepo.AddBook(ctx, book)
	logFatal(err)

	json.NewEncoder(w).Encode(bookID)
}

// UpdateBook updates existing book by id
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var book models.Book
	spew.Dump(json.NewDecoder(r.Body).Decode(&book), "JJJJJJJJ")
	json.NewDecoder(r.Body).Decode(&book)
	bookRepo := repository.BookRepository{}
	rowsUpdated, err := bookRepo.UpdateBook(ctx, book)
	logFatal(err)

	if book.Available {
		notifierapi.PubSub.Publish(book.ID, "Available")
	}

	json.NewEncoder(w).Encode(rowsUpdated)
}

// RemoveBook deletes book by id
func RemoveBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)

	bookRepo := repository.BookRepository{}
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	rowsDeleted, err := bookRepo.RemoveBook(ctx, id)
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
