package pkg

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/tousart/sso/config"
)

func ConnectToDB(cfg *config.Config) (*sql.DB, error) {
	log.Printf("USER NAME: %s", cfg.Postgres.User)
	address := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName)

	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %v", err)
	}

	return db, nil
}
