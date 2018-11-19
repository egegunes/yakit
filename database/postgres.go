package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func New(host string, name string, user string, password string) (*sql.DB, error) {
	auth := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, password, name, host)

	db, err := sql.Open("postgres", auth)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
