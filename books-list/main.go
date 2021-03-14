package main

import (
	"github.com/MeleshkoYuliya/golang/book-list/books-list"
	"github.com/MeleshkoYuliya/golang/book-list/driver"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.ConnectDB()
}

func main() {
	books.InitAPI()
}
