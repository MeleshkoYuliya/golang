package books

import (
	"context"
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

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type booksService struct {
	db          *sql.DB
	books       []models.Book
	pubSub      notifier.PubSub
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
	s.pubSub = notifier.NewPubSub()
	s.db = driver.GetDB()
	router := mux.NewRouter()
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books", AddBook).Methods("POST")
	router.HandleFunc("/books", UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", RemoveBook).Methods("DELETE")
	router.HandleFunc("/subscribers", CreateSubscriber).Methods("POST")
	router.HandleFunc("/suscriptions", SendNotification).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetBooks returns list of books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	r = r.WithContext(ctx)

	bookRepo := bookRepository.BookRepository{}
	books, err := bookRepo.GetBooks(ctx)
	logFatal(err)
	json.NewEncoder(w).Encode(books)
}

// GetBook returns one book by id
func GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	bookRepo := bookRepository.BookRepository{}
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
	bookRepo := bookRepository.BookRepository{}
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
	bookRepo := bookRepository.BookRepository{}
	rowsUpdated, err := bookRepo.UpdateBook(ctx, book)
	logFatal(err)

	if book.Available {
		s.pubSub.Publish(book.ID, "Available")
	}

	json.NewEncoder(w).Encode(rowsUpdated)
}

// RemoveBook deletes book by id
func RemoveBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)

	bookRepo := bookRepository.BookRepository{}
	id, err := strconv.Atoi(params["id"])
	logFatal(err)

	rowsDeleted, err := bookRepo.RemoveBook(ctx, id)
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}

// CreateSubscriber creates new book subscriber
func CreateSubscriber(w http.ResponseWriter, r *http.Request) {
	s.db = driver.GetDB()
	var subscriber models.Subscriber

	json.NewDecoder(r.Body).Decode(&subscriber)

	err := s.db.QueryRow("insert into public.subscribers (email, book_id) values($1, $2) RETURNING id;",
		subscriber.Email, subscriber.BookID).Scan(&subscriber.ID)

	logFatal(err)

	go func(email string) {
		bookCh := s.pubSub.Subscribe(subscriber.BookID)
		for b := range bookCh {
			callBackF(b, email, subscriber.BookID)
		}

	}(subscriber.Email)

	json.NewEncoder(w).Encode(subscriber.ID)
}

// SendNotification send notification on email for each subscriber
func SendNotification(w http.ResponseWriter, r *http.Request) {
	s.db = driver.GetDB()
	var subscriber models.Subscriber

	var bookID int
	json.NewDecoder(r.Body).Decode(&bookID)
	rows, err := s.db.Query("SELECT * from public.subscribers WHERE book_id=$1", bookID)
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&subscriber.ID, &subscriber.Email, &subscriber.BookID)
		logFatal(err)
		s.subscribers = append(s.subscribers, subscriber)
	}

	spew.Dump(s.subscribers, "s.subscribers")
	for _, sub := range s.subscribers {
		fmt.Printf("Отправлена нотификация на почту %v. Книга %v теперь доступна\n", sub.Email, sub.BookID)
	}

}

func callBackF(b interface{}, email string, bookID int) {
	fmt.Printf("Подписка на книгу %v по почте %v \n", bookID, email)
}
