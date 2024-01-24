package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	Firstname         string    `json:"firstname"`
	Lastname          string    `json:"lastname"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func newAccount(firstname, lastname, password string) (*Account, error) {
	encPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		Firstname:         firstname,
		Lastname:          lastname,
		Number:            rand.Int63n(100000),
		EncryptedPassword: string(encPw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
