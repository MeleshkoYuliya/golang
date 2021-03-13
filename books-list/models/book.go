package models

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Available bool   `json:"available"`
}
