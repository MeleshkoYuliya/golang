package bookRepository

import (
	"database/sql"
	"log"
	"main/models"
)

type BookRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (b BookRepository) GetBooks(db *sql.DB, book models.Book, books []models.Book) []models.Book {
	books = []models.Book{}

	rows, err := db.Query("SELECT * FROM public.books_list")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Available)
		logFatal(err)
		books = append(books, book)
	}

	return books
}

func (b BookRepository) GetBook(db *sql.DB, book models.Book, id int) models.Book {
	rows := db.QueryRow("SELECT * from public.books_list WHERE id=$1", id)
	err := rows.Scan(&book.ID, &book.Title, &book.Available)
	logFatal(err)

	return book
}

func (b BookRepository) AddBook(db *sql.DB, book models.Book) int {
	err := db.QueryRow("insert into public.books_list (title, available) values($1, $2) RETURNING id;",
		book.Title, book.Available).Scan(&book.ID)

	logFatal(err)
	return book.ID
}

func (b BookRepository) UpdateBook(db *sql.DB, book models.Book) int64 {
	result, err := db.Exec("UPDATE public.books_list set title=$1, available=$2 WHERE id=$3 RETURNING id",
		&book.Title, &book.Available, &book.ID)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}

func (b BookRepository) RemoveBook(db *sql.DB, id int) int64 {
	result, err := db.Exec("DELETE from public.books_list WHERE id=$1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted
}
