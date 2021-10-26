package repository

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/database/query"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
)

type UserRepositoryImpl struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewUserRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) (model.UserRepository, error) {
	var target UserRepositoryImpl
	if dbHandler == nil {
		return nil, errors.New("can't init user repository")
	}
	target.h = dbHandler
	target.l = log
	return &target, nil
}

func (ur *UserRepositoryImpl) Save(ctx context.Context, login string, pass string) error {
	if login == "" {
		ur.l.Info("UserRepository: empty login authorization attempt")
		return errors.New("can't register empty login")
	}
	if pass == "" {
		ur.l.Info("UserRepository: empty pass authorization  attempt")
		return errors.New("pass cannot be empty")
	}
	err := ur.h.Execute(ctx, query.CreateUser, login, pass)
	// TODO: нужен  userID в ответ
	// TODO: нужна транзакция
	userID := 0
	err = ur.h.Execute(ctx, query.CreateAccount, userID)

	return err
}

func (ur *UserRepositoryImpl) Check(ctx context.Context, login string, pass string) (bool, error) {
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
