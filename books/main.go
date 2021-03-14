package main

import (
	bookApi "github.com/MeleshkoYuliya/golang/books/book-api"
	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.GetDB()
}

func main() {
	bookApi.InitAPI()
}
