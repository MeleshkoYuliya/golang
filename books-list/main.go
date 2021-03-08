package main

import (
	"main/books-list"
	"main/driver"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.ConnectDB()
}

func main() {
	books.InitAPI()
}
