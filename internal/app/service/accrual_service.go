package service

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"go.uber.org/zap"
	"time"
)

type AccrualClient interface {
	GetAccrual(ctx context.Context, orderNum string) (*dto.Accrual, error)
}

type GophermartClient interface {
	ProcessRequest(ctx context.Context, orderNum string) bool
}

type AccrualService struct {
	dbOrder          model.OrderRepository
	dbBalance        model.BalanceRepository
	accrualClient    AccrualClient
	gophermartClient GophermartClient
	log              *infrastructure.Logger
}

func NewAccrualService(
	orderRepo model.OrderRepository,
	balanceRepo model.BalanceRepository,
	accrualClient AccrualClient,
	gophermartClient GophermartClient,
	log *infrastructure.Logger,
) *AccrualService {
	var target AccrualService
	target.dbOrder = orderRepo
	target.dbBalance = balanceRepo
	target.log = log
	target.accrualClient = accrualClient
	target.gophermartClient = gophermartClient
	return &target
}

func (s *AccrualService) StartProcessJob(latencyInSec time.Duration) {
	t := time.NewTicker(latencyInSec * time.Second)
	defer t.Stop()
	for true {
		_ = <-t.C
		ctx := context.Background()
		s.process(ctx)
	}
}

func (s *AccrualService) process(ctx context.Context) {
	// Выполнить обработку
	s.log.Debug("AccrualService: process. Start process job")
	orderList, err := s.dbOrder.FindNotProcessed(ctx)
	if err != nil {
		s.log.Error("AccrualService: process. Can't get order list", zap.Error(err))
		return
	}

	for _, order := range orderList {
		s.gophermartClient.ProcessRequest(ctx, order.Num)
	}

	s.log.Debug("AccrualService: process. Process job finished")
}

func (s *AccrualService) ProcessOrder(ctx context.Context, orderNum string) error {
	s.log.Debug("AccrualService: processOrder. Request")
	accrual, err := s.accrualClient.GetAccrual(ctx, orderNum)
	if err != nil {
		s.log.Error("AccrualService: processOrder. Can't get accruals from remote service", zap.Error(err))
		return err
	}
	order, err := s.dbOrder.LockOrder(ctx, orderNum)
	if err != nil {
		s.log.Error("AccrualService: processOrder. Can't lock order", zap.Error(err))
		return err
	}
	// Начисление делаем только для статуса Processed
	if accrual.Status == model.OrderStatusProcessed && order.Status != model.OrderStatusProcessed {
		account, err := s.dbBalance.LockAccount(ctx, order.UserID)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't lock account", zap.Error(err))
			return err
		}

		operation := model.Operation{
			AccountID:     account.ID,
			Amount:        accrual.Accrual,
			OrderID:       order.ID,
			OrderNum:      order.Num,
			OperationType: model.OperationCredit,
			ProcessedAt:   time.Now().Truncate(time.Second),
		}
		account.Balance += accrual.Accrual
		account.Credit += accrual.Accrual

		order.Status = accrual.Status
		order.UpdatedAt = time.Now().Truncate(time.Second)
		s.dbBalance.CreateOperation(ctx, &operation)
		s.dbBalance.SaveAccount(ctx, account)
		s.dbOrder.Save(ctx, order)
	} else if accrual.Status == model.OrderStatusProcessing || accrual.Status == model.OrderStatusRegistered || accrual.Status == model.OrderStatusInvalid {
		order.Status = accrual.Status
		order.UpdatedAt = time.Now().Truncate(time.Second)
		s.dbOrder.Save(ctx, order)
	} else {
		s.log.Error("AccrualService: processOrder. Recieved unexpected status", zap.String("OrderNum", order.Num), zap.String("Status", accrual.Status))
		return errors.New("recieved unexpected status")
	}
	s.log.Debug("AccrualService: processOrder. Success")
	return nil
}