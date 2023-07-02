package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)
type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password string `json:"password"`
}
type Account struct {
		ID                int       `json:"id"`
		FirstName         string    `json:"firstName"`
		LastName          string    `json:"lastName"`
		Number            int64     `json:"number"`
		EncryptedPassword string    `json:"-"`
		Balance           int64     `json:"balance"`
		CreatedAt         time.Time `json:"createdAt"`
	}
func (a *Account) ValidPassword(pw string) bool {
		return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
	}
type LoginRequest struct {
		Number int64 `json:"number"`
		Password string  `json:"password"`
	}
type LoginResponse struct {
	Token string `json:"token"`
	Number string  `json:"number"`
	}
type TransferRequest struct {
		ToAccount int `json:"toAccount"`
		Amount    int `json:"amount"`
	}

func NewAccount(FirstName , lastName ,password string)(*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if  err !=nil {
		return nil ,err
	}
	return &Account{
        FirstName: FirstName,
        LastName:  lastName,
        Number:    int64(rand.Intn(100000000)),
		EncryptedPassword: string(encpw),
		CreatedAt: time.Now().UTC(),
    },nil
}