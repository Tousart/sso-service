package models

type User struct {
	UserID       string
	Login        string
	Email        string
	HashPassword []byte
}
