package main

import (
	"github.com/MeleshkoYuliya/golang/common/driver"
	"github.com/MeleshkoYuliya/golang/notifier/notifierapi"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	driver.GetDB()
}

func main() {
	notifierapi.InitAPI()
}
