package main

import (
	"math/rand"
	"time"
)

type TransereRequest struct {
	To     int `json:"to"`
	Amount uint32 `json:"amount"` 
}


type createAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       uint8    `json:"age"`
}

type Account struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Age        uint8     `json:"age"`
	Balance    uint32    `json:"balance"`
	BankNumber int       `json:"bank_number"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewAccount(firstName string, lastName string, age uint8) *Account {

	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Balance:   0,
		BankNumber: rand.Intn(10000),
		CreatedAt: time.Now().UTC(),
	}

}
