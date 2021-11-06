package dto

import (
	"errors"
	"time"
)

type Error struct {
	Msg string `json:"msg"`
}

type User struct {
	ID    int
	Login string `json:"login"`
	Pass  string `json:"password"`
}

type Order struct {
	Num      string    `json:"number"`
	UserID   int       `json:"-"`
	Status   string    `json:"status"`
	Accrual  int       `json:"accrual"`
	UploadAt time.Time `json:"upload_at"`
}

/*
 {
          "order": "2377225624",
          "sum": 500,
          "status": "PROCESSED",
          "processed_at": "2020-12-09T16:09:57+03:00"
      }
*/
type Withdrawal struct {
	OrderNum    string    `json:"order"`
	Amount      int       `json:"sum"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Withdraw struct {
	OrderNum string `json:"order"`
	Amount   int    `json:"sum"`
}

type Balance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNotFound = errors.New("no rows in result set")
var ErrBadParam = errors.New("bad param occured")
var ErrUnauthorized = errors.New("User unauthorized")

var ErrOrderRegistered = errors.New("order registered early")
var ErrOrderRegisteredByAnotherUser = errors.New("order registered early by another user")

var ErrNotEnoughFunds = errors.New("not enougth founds")
var ErrBadOrderNum = errors.New("bad order num")
