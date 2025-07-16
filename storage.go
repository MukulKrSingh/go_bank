package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(id int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
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
func (s *PostgresStore) Init() error {
	return s.createTableAccount()

}

func (s *PostgresStore) createTableAccount() error {
	query := `CREATE TABLE IF NOT EXISTS account(
	id SERIAL PRIMARY KEY,
	first_name varchar(50),
	last_name varchar(50),
	number SERIAL,
	balance SERIAL,
	created_at timestamp
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}

	for rows.Next() {
		acc := new(Account)
		if err := rows.Scan(&acc.Id,
			&acc.FirstName,
			&acc.LastName,
			&acc.Number,
			&acc.Balance,
			&acc.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (first_name , last_name ,number, balance , created_at)
	VALUES ($1, $2 , $3 , $4, $5)
	`
	resp, err := s.db.Query(
		query,
		acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}
func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
