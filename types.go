package main

import (
	"math/rand"
)

type account struct {
	ID      int    `json:"id"`
	Number  int64  `json:"number"`
	Email   string `json:"email"`
	Balance int64  `json:"balance"`
}

func NewAccount(Email string) *account {
	return &account{
		ID:     rand.Intn(10000),
		Email:  Email,
		Number: rand.Int63n(10000),
	}
}
