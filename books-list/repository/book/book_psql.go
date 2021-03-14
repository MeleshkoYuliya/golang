package bookRepository

import (
	"context"
	"database/sql"
	"log"

	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/MeleshkoYuliya/golang/common/models"
)

type BookRepository struct {
	db *sql.DB
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b BookRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	var book models.Book
	books := []models.Book{}
	b.db = driver.GetDB()

	rows, err := b.db.QueryContext(ctx, "SELECT * FROM public.books_list")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Available)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (b BookRepository) GetBook(ctx context.Context, id int) (models.Book, error) {
	var book models.Book
	b.db = driver.GetDB()

	rows := b.db.QueryRowContext(ctx, "SELECT * from public.books_list WHERE id=$1", id)
	err := rows.Scan(&book.ID, &book.Title, &book.Available)
	if err != nil {
		return book, err
	}

	return book, nil
}

func (b BookRepository) AddBook(ctx context.Context, book models.Book) (int, error) {
	b.db = driver.GetDB()
	err := b.db.QueryRowContext(ctx, "insert into public.books_list (title, available) values($1, $2) RETURNING id;",
		book.Title, book.Available).Scan(&book.ID)

	if err != nil {
		return book.ID, err
	}

	return book.ID, nil
}

func (b BookRepository) UpdateBook(ctx context.Context, book models.Book) (int64, error) {
	b.db = driver.GetDB()
	result, err := b.db.ExecContext(ctx, "UPDATE public.books_list set title=$1, available=$2 WHERE id=$3 RETURNING id",
		&book.Title, &book.Available, &book.ID)

	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return rowsUpdated, err
	}

	return rowsUpdated, nil
}

func (b BookRepository) RemoveBook(ctx context.Context, id int) (int64, error) {
	b.db = driver.GetDB()
	result, err := b.db.ExecContext(ctx, "DELETE from public.books_list WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return rowsDeleted, err
	}

	return rowsDeleted, nil
}
