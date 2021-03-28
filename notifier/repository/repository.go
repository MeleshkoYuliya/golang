package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/MeleshkoYuliya/golang/common/models"
)

type NotifierRepository struct {
	db *sql.DB
}

func (n NotifierRepository) CreateSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	n.db = driver.GetDB()

	err := n.db.QueryRowContext(ctx, "insert into public.subscribers (email, book_id) values($1, $2) RETURNING id;",
		subscriber.Email, subscriber.BookID).Scan(&subscriber.ID)

	if err != nil {
		return err
	}

	fmt.Printf("Subscriber was added with id %v", subscriber.ID)
	return nil
}
