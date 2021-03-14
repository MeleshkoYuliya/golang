package main

import (
	books "github.com/MeleshkoYuliya/golang/book-list/book-api"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.ConnectDB()
}

func main() {
	books.InitAPI()
}
