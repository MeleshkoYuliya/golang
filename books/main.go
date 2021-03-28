package main

import (
	"github.com/MeleshkoYuliya/golang/books/bookapi"
	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	driver.GetDB()
	bookapi.InitAPI()
}
