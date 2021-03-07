// package books

// import (
// 	"database/sql"
// 	"io/ioutil"
// 	"log"
// 	"main/books-list/driver"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// )

// func NewMock() (*sql.DB, sqlmock.Sqlmock) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	return db, mock
// }

// func TestGetBooks(t *testing.T) {
// 	db = driver.ConnectDB()
// 	req, err := http.NewRequest("GET", "localhost:8080", nil)
// 	if err != nil {
// 		t.Fatalf("could not created request: %v", err)
// 	}

// 	println(req)
// 	rec := httptest.NewRecorder()
// 	println(rec)
// 	GetBooks(rec, req)

// 	res := rec.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Errorf("expected status OK; got %v", res.Status)
// 	}

// 	_, err = ioutil.ReadAll(res.Body)

// 	if err != nil {
// 		t.Fatalf("could not read response: %v", err)
// 	}

// }
