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
	enable           bool
}

func NewAccrualService(
	orderRepo model.OrderRepository,
	balanceRepo model.BalanceRepository,
	accrualClient AccrualClient,
	gophermartClient GophermartClient,
	log *infrastructure.Logger,
	enable bool,
) *AccrualService {
	var target AccrualService
	target.dbOrder = orderRepo
	target.dbBalance = balanceRepo
	target.log = log
	target.accrualClient = accrualClient
	target.gophermartClient = gophermartClient
	target.enable = enable
	return &target
}

func (s *AccrualService) StartProcessJob(latency time.Duration) {
	if !s.enable {
		return
	}
	t := time.NewTicker(latency * time.Second)
	defer t.Stop()
	for {
		<-t.C
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
	accrual, err := s.accrualClient.GetAccrual(context.Background(), orderNum)
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
		err = s.dbBalance.CreateOperation(ctx, &operation)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't create operation", zap.Error(err))
			return err
		}
		err = s.dbBalance.SaveAccount(ctx, account)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't save account", zap.Error(err))
			return err
		}
		err = s.dbOrder.UpdateStatus(ctx, order)
		if err != nil {
			s.log.Error("AccrualService: processOrder. Can't save order", zap.Error(err))
			return err
		}
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
