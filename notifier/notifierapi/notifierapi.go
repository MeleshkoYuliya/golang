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
	"github.com/MeleshkoYuliya/golang/notifier/repository"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

type notifierService struct {
	db    *sql.DB
	books []models.Book
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
	ctx := r.Context()
	var subscriber models.Subscriber
	json.NewDecoder(r.Body).Decode(&subscriber)

	notifierRepo := repository.NotifierRepository{}

	err := notifierRepo.CreateSubscriber(ctx, subscriber)

	if err != nil {
		fmt.Printf("Failed to create subscriber. Server responded with status %v", http.StatusBadRequest)
	}

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
	ctx := r.Context()
	var bookID int
	json.NewDecoder(r.Body).Decode(&bookID)

	notifierRepo := repository.NotifierRepository{}

	subscribers, err := notifierRepo.GetSubscribersByBookID(ctx, bookID)
	if err != nil {
		fmt.Printf("Failed to get subscribers by bookID %v. Server responded with status %v", bookID, http.StatusBadRequest)
	}

	for _, sub := range subscribers {
		fmt.Printf("Отправлена нотификация на почту %v. Книга %v теперь доступна\n", sub.Email, sub.BookID)
	}

}

func callBackF(b interface{}, email string, bookID int) {
	fmt.Printf("Подписка на книгу %v по почте %v \n", bookID, email)
}
