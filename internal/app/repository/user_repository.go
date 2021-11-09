package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
	"go.uber.org/zap"
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

func (ur *UserRepositoryImpl) Save(ctx context.Context, login string, pass string) (int, error) {
	var userID int
	row, err := ur.h.QueryRow(ctx, GetNextUserID)
	if err != nil {
		ur.l.Error("UserRepository: cannt get userID")
		return 0, err
	}
	err = row.Scan(&userID)
	if err != nil {
		ur.l.Error("UserRepository: cannt get userID")
		return 0, err
	}

	err = ur.h.Execute(ctx, CreateUser, userID, login, pass)
	if err != nil {
		ur.l.Error("UserRepository: cannt create user", zap.Error(err))
		return 0, err
	}
	err = ur.h.Execute(ctx, CreateAccount, userID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return 0, &model.UniqueViolation
		}
	}
	if err != nil {
		ur.l.Error("UserRepository: cannt create account for user", zap.Error(err))
		return 0, err
	}
	return userID, nil
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
	row, err := ur.h.QueryRow(ctx, CheckUser, login, pass)
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

func (ur *UserRepositoryImpl) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {

	row, err := ur.h.QueryRow(ctx, GetUserByLogin, login)
	if err != nil {
		return nil, err
	}
	var res model.User
	err = row.Scan(&res.ID, &res.Login, &res.Pass)
	if err != nil && err.Error() == "no rows in result set" {
		return nil, &model.NoRowFound
	}

	if err != nil {
		return nil, nil
	}
	return &res, nil
}
