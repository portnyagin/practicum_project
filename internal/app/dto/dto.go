package dto

import "errors"

type Error struct {
	Msg string `json:"msg"`
}

type User struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}

var ErrDuplicateKey = errors.New("duplicate key")
var ErrNotFound = errors.New("no rows in result set")
var ErrBadParam = errors.New("bad param occured")
