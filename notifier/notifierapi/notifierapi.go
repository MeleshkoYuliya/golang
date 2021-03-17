package notifierapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/MeleshkoYuliya/golang/common/models"
	"github.com/MeleshkoYuliya/golang/notifier/pubsub"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type notifierService struct {
	db          *sql.DB
	books       []models.Book
	subscribers []models.Subscriber
}

var n notifierService
var PubSub pubsub.PubSub

func logFatal(err error) {
	if err != nil {
		spew.Dump(err)
	}
}

// InitAPI initiates routes
func InitAPI() {
	PubSub = pubsub.NewPubSub()
	n.db = driver.GetDB()
	router := mux.NewRouter()
	router.HandleFunc("/subscribers", CreateSubscriber).Methods("POST")
	router.HandleFunc("/suscriptions", SendNotification).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// CreateSubscriber creates new book subscriber
func CreateSubscriber(w http.ResponseWriter, r *http.Request) {
	n.db = driver.GetDB()
	var subscriber models.Subscriber

	json.NewDecoder(r.Body).Decode(&subscriber)

	err := n.db.QueryRow("insert into public.subscribers (email, book_id) values($1, $2) RETURNING id;",
		subscriber.Email, subscriber.BookID).Scan(&subscriber.ID)

	logFatal(err)

	go func(email string) {
		bookCh := PubSub.Subscribe(subscriber.BookID)
		for b := range bookCh {
			callBackF(b, email, subscriber.BookID)
		}

	}(subscriber.Email)

	json.NewEncoder(w).Encode(subscriber.ID)
}

// SendNotification send notification on email for each subscriber
func SendNotification(w http.ResponseWriter, r *http.Request) {
	n.db = driver.GetDB()
	var subscriber models.Subscriber

	var bookID int
	json.NewDecoder(r.Body).Decode(&bookID)
	rows, err := n.db.Query("SELECT * from public.subscribers WHERE book_id=$1", bookID)
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&subscriber.ID, &subscriber.Email, &subscriber.BookID)
		logFatal(err)
		n.subscribers = append(n.subscribers, subscriber)
	}

	for _, sub := range n.subscribers {
		fmt.Printf("Отправлена нотификация на почту %v. Книга %v теперь доступна\n", sub.Email, sub.BookID)
	}

}

func callBackF(b interface{}, email string, bookID int) {
	fmt.Printf("Подписка на книгу %v по почте %v \n", bookID, email)
}
