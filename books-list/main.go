package main

import (
	"main/books-list"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	books.InitAPI()
}
