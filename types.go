package main

import (
	"math/rand"
	"time"
)
type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`}
type Account struct {
		ID                int       `json:"id"`
		FirstName         string    `json:"firstName"`
		LastName          string    `json:"lastName"`
		Number            int64     `json:"number"`
		EncryptedPassword string    `json:"-"`
		Balance           int64     `json:"balance"`
		CreatedAt         time.Time `json:"createdAt"`
	}
func NewAccount(FirstName string, lastName string)(*Account, error) {
	return &Account{
        FirstName: FirstName,
        LastName:  lastName,
        Number:    int64(rand.Intn(100000000)),
		CreatedAt: time.Now().UTC(),
    },nil
}