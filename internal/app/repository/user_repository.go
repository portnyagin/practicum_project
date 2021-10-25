package repository

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/database/query"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
)

type UserRepository struct {
	h basedbhandler.DBHandler
}

func NewUserRepository(dbHandler basedbhandler.DBHandler) *UserRepository {
	var target UserRepository
	target.h = dbHandler
	return &target
}

func (ur *UserRepository) Save(ctx context.Context, login string, pass string) error {
	if login == "" {
		return errors.New("can't register empty login")
	}
	if pass == "" {
		return errors.New("pass cannot be empty")
	}
	err := ur.h.Execute(ctx, query.CreateUser, login, pass)
	return err
}

func (ur *UserRepository) Check(ctx context.Context, login string, pass string) (bool, error) {
	if login == "" {
		return false, errors.New("can't register empty login")
	}
	if pass == "" {
		return false, errors.New("pass cannot be empty")
	}
	row, err := ur.h.QueryRow(ctx, query.CheckUser, login, pass)
	if err != nil {
		return false, err
	}

	var res int
	err = row.Scan(&res)
	if err != nil {
		return false, nil
	}
	return true, nil
}
