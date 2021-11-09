package repository

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
	"go.uber.org/zap"
)

type BalanceRepository struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewBalanceRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) (model.BalanceRepository, error) {
	var target BalanceRepository
	if dbHandler == nil {
		return nil, errors.New("can't init balance repository")
	}
	target.h = dbHandler
	target.l = log
	return &target, nil
}

func (r *BalanceRepository) FindWithdrawalByUser(ctx context.Context, userID int) ([]model.Withdrawal, error) {
	rows, err := r.h.Query(ctx, GetWithdrawalByUser, userID)
	var resArray []model.Withdrawal
	if err != nil {
		r.l.Error("BalanceRepository: request error", zap.String("query", GetWithdrawalByUser), zap.Int("userID", userID), zap.Error(err))
		return nil, err
	}
	for rows.Next() {
		var o model.Withdrawal
		err := rows.Scan(&o.OrderNum, &o.Amount, &o.Status, &o.ProcessedAt)
		if err != nil {
			r.l.Error("BalanceRepository: scan rows error", zap.String("query", GetWithdrawalByUser), zap.Int("userID", userID), zap.Error(err))
			break
		}
		resArray = append(resArray, o)
	}
	return resArray, nil
}

func (r *BalanceRepository) LockAccount(ctx context.Context, userID int) (*model.Account, error) {
	row, err := r.h.QueryRow(ctx, GetAccountForUpdate, userID)
	if err != nil {
		r.l.Error("BalanceRepository: cannt get account for update", zap.Error(err))
		return nil, err
	}
	account := model.Account{}
	err = row.Scan(&account.ID, &account.UserID, &account.Balance, &account.Debit, &account.Credit)
	if err != nil {
		r.l.Error("BalanceRepository: cannt get account for update", zap.Error(err))
		if err.Error() == "no rows in result set" {
			return nil, &model.NoRowFound
		} else {
			return nil, err
		}
	}

	return &account, nil
}

func (r *BalanceRepository) SaveAccount(ctx context.Context, account *model.Account) error {
	err := r.h.Execute(ctx, UpdateAccountForUser, account.UserID, account.Balance, account.Debit, account.Credit)
	if err != nil {
		r.l.Error("UserRepository: cannt create user", zap.Error(err))
		return err
	}
	return nil
}

func (r *BalanceRepository) CreateOperation(ctx context.Context, operation *model.Operation) error {
	err := r.h.Execute(ctx, CreateOperation,
		operation.AccountID,
		operation.OrderID,
		operation.OrderNum,
		operation.OperationType,
		operation.Amount,
		operation.ProcessedAt)
	if err != nil {
		r.l.Error("BalanceRepository: cannt create operation", zap.Error(err))
		return err
	}
	return nil
}

func (r *BalanceRepository) GetAccount(ctx context.Context, userID int) (*model.Account, error) {
	row, err := r.h.QueryRow(ctx, GetAccount, userID)
	if err != nil {
		r.l.Error("BalanceRepository: cannt get account for update", zap.Error(err))
		return nil, err
	}
	account := model.Account{}
	err = row.Scan(&account.ID, &account.UserID, &account.Balance, &account.Debit, &account.Credit)
	if err != nil {
		r.l.Error("BalanceRepository: cannt get account for update", zap.Error(err))
		if err.Error() == "no rows in result set" {
			return nil, &model.NoRowFound
		} else {
			return nil, err
		}
	}

	return &account, nil
}
