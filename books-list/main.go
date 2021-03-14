package main

import (
	books "github.com/MeleshkoYuliya/golang/books-list/book-api"
	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.GetDB()
}

func main() {
	books.InitAPI()
}
