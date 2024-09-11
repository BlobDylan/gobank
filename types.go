package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type loginrequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type transferRequest struct {
	To     int   `json:"to"`
	Amount int64 `json:"amount"`
}

type createAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type account struct {
	ID           int       `json:"id"`
	Number       int64     `json:"number"`
	Email        string    `json:"email"`
	EncryptedPwd string    `json:"-"`
	Balance      int64     `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewAccount(Email, Password string) (*account, error) {
	EncryptedPwd, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &account{
		Email:        Email,
		EncryptedPwd: string(EncryptedPwd),
		Number:       rand.Int63n(10000),
		CreatedAt:    time.Now().UTC(),
	}, nil
}
