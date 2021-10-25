package model

import "context"

type Repository interface {
	Save(ctx context.Context, login string, pass string) error
	Check(ctx context.Context, login string, pass string) (bool, error)
}

type User struct {
	Login  string
	Pass   string
	Active bool
}
