package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DB struct {
	host     string
	name     string
	user     string
	password string
}

func New(host string, name string, user string, password string) *DB {
	return &DB{
		host:     host,
		name:     name,
		user:     user,
		password: password,
	}
}

func (d DB) Open() (*sql.DB, error) {
	authString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", d.user, d.password, d.name, d.host)

	db, err := sql.Open("postgres", authString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
