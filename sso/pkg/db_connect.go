package pkg

import (
	"database/sql"
	"fmt"
)

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://user:password@postgres:5432/auth_psql?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %v", err)
	}

	return db, nil
}
