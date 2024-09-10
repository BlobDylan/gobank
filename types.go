package main

import (
	"math/rand"
	"time"
)

type transferRequest struct {
	To     int   `json:"to"`
	Amount int64 `json:"amount"`
}

type createAccountRequest struct {
	Email string `json:"email"`
}

type account struct {
	ID        int       `json:"id"`
	Number    int64     `json:"number"`
	Email     string    `json:"email"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(Email string) *account {
	return &account{
		Email:     Email,
		Number:    rand.Int63n(10000),
		CreatedAt: time.Now().UTC(),
	}
}
