package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccoutn(id int) error
	UpdateAccoutn(*Account) error
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to open db : %v", err)
		return nil, err

	}

	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to open db : %v", err)
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}
