package books

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	bookRepository "main/books-list/repository/book"
	"main/driver"
	"main/models"
	"main/notifier"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var books []models.Book
var subscribers []models.Subscriber
var pubSub notifier.PubSub

type booksService struct {
	db *sql.DB
}

var s booksService

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// InitAPI initiates routes
func InitAPI() {
	pubSub = notifier.NewPubSub()
	s.db = driver.GetDB()
	router := mux.NewRouter()
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books", AddBook).Methods("POST")
	router.HandleFunc("/books", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", RemoveBook).Methods("DELETE")
	router.HandleFunc("/subscribers", CreateSubscriber).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetBooks returns list of books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	books = []models.Book{}
	bookRepo := bookRepository.BookRepository{}
	books = bookRepo.GetBooks(s.db, book, books)

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
	book = bookRepo.GetBook(s.db, book, id)

	json.NewEncoder(w).Encode(book)
}

// AddBook adds new book to books list
func AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)
	bookRepo := bookRepository.BookRepository{}
	bookID = bookRepo.AddBook(s.db, book)

	json.NewEncoder(w).Encode(bookID)
}

// UpdateBook updates existing book by id
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	bookRepo := bookRepository.BookRepository{}
	rowsUpdated := bookRepo.UpdateBook(s.db, book)
	if book.Available {
		pubSub.Publish(book.ID, "Available")
	}

	json.NewEncoder(w).Encode(rowsUpdated)
}

// RemoveBook deletes book by id
func RemoveBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	bookRepo := bookRepository.BookRepository{}
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	rowsDeleted := bookRepo.RemoveBook(s.db, id)

	json.NewEncoder(w).Encode(rowsDeleted)
}

// CreateSubscriber creates new book subscriber
func CreateSubscriber(w http.ResponseWriter, r *http.Request) {
	db := driver.GetDB()
	var subscriber models.Subscriber

	json.NewDecoder(r.Body).Decode(&subscriber)

	err := db.QueryRow("insert into public.subscribers (email, book_id) values($1, $2) RETURNING id;",
		subscriber.Email, subscriber.BookID).Scan(&subscriber.ID)

	logFatal(err)

	go func(email string) {
		bookCh := pubSub.Subscribe(subscriber.BookID)
		for b := range bookCh {
			callBackF(b, email, subscriber.BookID)
		}

	}(subscriber.Email)

	json.NewEncoder(w).Encode(subscriber.ID)
}

func callBackF(b interface{}, email string, bookID int) {
	fmt.Printf("Подписка на книгу %v по почте %v", bookID, email)
}
