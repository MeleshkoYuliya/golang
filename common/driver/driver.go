package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

var db *sql.DB

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init() {
	ConnectDB()
}

func ConnectDB() *sql.DB {
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		fmt.Printf("Failed to connect db")
	}

	db, err = sql.Open("postgres", pgUrl)
	if err != nil {
		fmt.Printf("Failed to get db")
	}

	db.SetMaxOpenConns(2)
	err = db.Ping()
	logFatal(err)
	return db
}

func GetDB() *sql.DB {
	if db == nil {
		return ConnectDB()
	}
	return db
}
