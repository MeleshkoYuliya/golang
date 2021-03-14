package main

import (
	"github.com/MeleshkoYuliya/golang/books/bookapi"
	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.GetDB()
}

func main() {
	bookapi.InitAPI()
}
