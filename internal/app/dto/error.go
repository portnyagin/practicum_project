package dto

import (
	"errors"
)

type Error struct {
	Msg string `json:"msg"`
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNotFound = errors.New("no rows in result set")
var ErrBadParam = errors.New("bad param occured")
var ErrUnauthorized = errors.New("User unauthorized")

var ErrOrderRegistered = errors.New("order registered early")
var ErrOrderRegisteredByAnotherUser = errors.New("order registered early by another user")

var ErrNotEnoughFunds = errors.New("not enougth founds")
var ErrBadOrderNum = errors.New("bad order num")

var ErrTooManyRequest = errors.New("too many request to remote service")
var ErrRemoteServiceError = errors.New("remote service error")
