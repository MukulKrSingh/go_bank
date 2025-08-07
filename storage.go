package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(id int) error
	UpdateAccount(ctx context.Context, tx *sql.Tx, account *Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(number int64) (*Account, error)
	HandleTransfer(context.Context, *Account, *Account, *TransferRequest) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://postgres:postgresd@localhost:5432/postgres?sslmode=disable"

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
	encrypted_password varchar(100),
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

func (s *PostgresStore) HandleTransfer(ctx context.Context, sen *Account, rec *Account, tr *TransferRequest) error {

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Fetch sender account with lock
	fromAccount := sen
	req := tr
	// Fetch receiver account with lock
	toAccount := rec
	// Check sufficient funds
	if (fromAccount.Balance) < int64(req.Amount) {
		return fmt.Errorf("insufficient funds")
	}

	// Adjust balances in-memory
	fromAccount.Balance -= int64(req.Amount)
	toAccount.Balance += int64(req.Amount)

	// Update accounts with new balances
	if err := s.UpdateAccount(ctx, tx, fromAccount); err != nil {
		return err
	}
	if err := s.UpdateAccount(ctx, tx, toAccount); err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()

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
		acc, err = scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (
	first_name , 
	last_name , 
	encrypted_password,
	number, 
	balance,
	created_at
	)
	VALUES ($1, $2 , $3 , $4, $5, $6)
	`
	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.EncryptedPassword,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

// Assume Account struct has: ID int64, Balance int64 (in cents, recommended)
func (s *PostgresStore) UpdateAccount(ctx context.Context, tx *sql.Tx, account *Account) error {
	query := "UPDATE account SET balance = $1 WHERE id = $2"
	_, err := tx.ExecContext(ctx, query, account.Balance, account.Id)
	return err
}
func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM account			
			WHERE  id = $1`
	resp, err := s.db.Query(query, id)
	if err != nil {
		return err
	}
	log.Printf("%+v", resp)
	return nil
}
func (s *PostgresStore) GetAccountById(id int) (*Account, error) {

	query := `SELECT * FROM account WHERE id = $1`

	rows := s.db.QueryRow(query, id)

	var acc Account
	err := rows.Scan(
		&acc.Id,
		&acc.FirstName,
		&acc.LastName,
		&acc.EncryptedPassword,
		&acc.Number,
		&acc.Balance,
		&acc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (s *PostgresStore) GetAccountByNumber(number int64) (*Account, error) {
	query := `SELECT * FROM account WHERE number = $1`

	rows := s.db.QueryRow(query, number)

	var acc Account
	err := rows.Scan(
		&acc.Id,
		&acc.FirstName,
		&acc.LastName,
		&acc.EncryptedPassword,
		&acc.Number,
		&acc.Balance,
		&acc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &acc, nil

}
func scanIntoAccount(rows *sql.Rows) (*Account, error) {

	account := new(Account)

	err := rows.Scan(
		&account.Id,
		&account.FirstName,
		&account.LastName,
		&account.EncryptedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
