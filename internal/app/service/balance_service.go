package service

import (
	"context"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"go.uber.org/zap"
	"time"
)

type BalanceService struct {
	dbBalance model.BalanceRepository
	log       *infrastructure.Logger
}

func NewBalanceService(balanceRepo model.BalanceRepository, log *infrastructure.Logger) *BalanceService {
	var target BalanceService
	target.dbBalance = balanceRepo
	target.log = log
	return &target
}

func (s *BalanceService) mapWithdrawalModelToDTO(src model.Withdrawal) dto.Withdrawal {
	return dto.Withdrawal{
		OrderNum:    src.OrderNum,
		Amount:      src.Amount,
		Status:      src.Status,
		ProcessedAt: src.ProcessedAt,
	}
}

func (s *BalanceService) mapWithdrawalListModelToDTO(src []model.Withdrawal) (resList []dto.Withdrawal) {
	for _, o := range src {
		resList = append(resList, s.mapWithdrawalModelToDTO(o))
	}
	return resList
}

func (s *BalanceService) GetCurrentBalance(ctx context.Context, userID int) (*dto.Balance, error) {
	if userID == 0 {
		s.log.Debug("BalanceService: GetCurrentBalance. got nil userID")
		return nil, dto.ErrBadParam
	}

	account, err := s.dbBalance.GetAccount(ctx, userID)
	if err != nil {
		s.log.Debug("BalanceService: GetCurrentBalance. Cannt get current balance")
		return nil, err
	}

	return &dto.Balance{
		Current:   account.Balance,
		Withdrawn: account.Credit,
	}, nil

}

func (s *BalanceService) Withdraw(ctx context.Context, obj *dto.Withdraw, userID int) error {
	if userID == 0 {
		s.log.Debug("BalanceService: Withdraw. got nil userID")
		return dto.ErrBadParam
	}
	if obj == nil {
		s.log.Debug("BalanceService: Withdraw. got nil order")
		return dto.ErrBadParam
	}
	if !CheckOrderNum(obj.OrderNum) {
		s.log.Debug("BalanceService: Withdraw. Order num validation error", zap.String("orderNum", obj.OrderNum))
		return dto.ErrBadOrderNum
	}
	account, err := s.dbBalance.LockAccount(ctx, userID)
	if err != nil {
		s.log.Error("BalanceService: Withdraw. Unexpected error", zap.Error(err))
		return err
	}

	if account.Balance < obj.Amount {
		s.log.Debug("BalanceService: Withdraw. In account not enough funds")
		return dto.ErrNotEnoughFunds
	}

	operation := model.Operation{
		AccountID:     account.ID,
		Amount:        obj.Amount,
		OrderNum:      obj.OrderNum,
		OperationType: model.OperationDebit,
		ProcessedAt:   time.Now().Truncate(time.Second),
	}
	err = s.dbBalance.CreateOperation(ctx, &operation)

	if err != nil {
		s.log.Error("BalanceService: Withdraw. Can't save operation", zap.Error(err))
		return err
	}
	account.Balance -= obj.Amount
	account.Debit += obj.Amount
	err = s.dbBalance.SaveAccount(ctx, account)
	if err != nil {
		s.log.Error("BalanceService: Withdraw. Can't save account", zap.Error(err))
		return err
	}

	return nil
}

func (s *BalanceService) GetWithdrawalsList(ctx context.Context, userID int) ([]dto.Withdrawal, error) {
	if userID == 0 {
		s.log.Debug("BalanceService: GetWithdrawalsList. got nil userID")
		return nil, dto.ErrBadParam
	}

	withdrawalList, err := s.dbBalance.FindWithdrawalByUser(ctx, userID)
	if err != nil {
		s.log.Error("BalanceService: GetWithdrawalsList. Can't get withdrawal list",
			zap.Int("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}

	resList := s.mapWithdrawalListModelToDTO(withdrawalList)
	return resList, nil
}
