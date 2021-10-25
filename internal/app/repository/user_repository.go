package repository

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/database/query"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
)

type UserRepository struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewUserRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) *UserRepository {
	var target UserRepository
	target.h = dbHandler
	target.l = log
	return &target
}

func (ur *UserRepository) Save(ctx context.Context, login string, pass string) error {
	if login == "" {
		ur.l.Info("UserRepository: empty login authorization attempt")
		return errors.New("can't register empty login")
	}
	if pass == "" {
		ur.l.Info("UserRepository: empty pass authorization  attempt")
		return errors.New("pass cannot be empty")
	}
	err := ur.h.Execute(ctx, query.CreateUser, login, pass)
	return err
}

func (ur *UserRepository) Check(ctx context.Context, login string, pass string) (bool, error) {
	if login == "" {
		ur.l.Info("UserRepository: empty login authorization attempt")
		return false, errors.New("can't register empty login")
	}
	if pass == "" {
		ur.l.Info("UserRepository: empty pass authorization  attempt")
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
