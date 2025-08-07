package main

import (
	"log"
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToAccount int `json:"to_account"`
	Amount    int `json:"amount"`
}
type CreateAccountRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Balance int64 `json:"balance,omitempty"`
	Password  string `json:"password"`
}
type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type Account struct {
	Id                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	EncryptedPassword string    `json:"-"`
	Number            int64     `json:"number"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(balance int64,firstName, lastName, password string) *Account {
	acc := Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1000000)),
		CreatedAt:         time.Now().UTC(),
		EncryptedPassword: password,
		Balance: balance,
	}
	log.Printf("Accont Created %+v", acc)
	return &acc
}
