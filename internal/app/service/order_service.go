package service

import (
	"context"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"go.uber.org/zap"
	"time"
)

type OrderService struct {
	dbOrder          model.OrderRepository
	log              *infrastructure.Logger
	EnableValidation bool
}

func NewOrderService(orderRepo model.OrderRepository, log *infrastructure.Logger, enableValidation bool) *OrderService {
	var target OrderService
	target.dbOrder = orderRepo
	target.log = log
	target.EnableValidation = enableValidation
	return &target
}

func (s *OrderService) mapOrderDTOtoModel(src *dto.Order) *model.Order {
	return &model.Order{
		UserID:   src.UserID,
		Num:      src.Num,
		Status:   src.Status,
		UploadAt: src.UploadAt,
	}
}

func (s *OrderService) mapOrderModeltoDTO(src *model.Order) *dto.Order {
	return &dto.Order{
		UserID:   src.UserID,
		Num:      src.Num,
		Status:   src.Status,
		Accrual:  src.Accrual,
		UploadAt: src.UploadAt.Truncate(time.Second),
	}
}

func (s *OrderService) mapOrderListModelToDTO(src []model.Order) (resList []dto.Order) {
	for _, o := range src {
		resList = append(resList, *s.mapOrderModeltoDTO(&o))
	}
	return resList
}

func (s *OrderService) Save(ctx context.Context, order *dto.Order) error {
	if order == nil {
		s.log.Debug("OrderService: Save. got nil order")
		return dto.ErrBadParam
	}
	if (order.UserID == 0) || (order.Num == "") {
		s.log.Debug("OrderService: Save. Validation error")
		return dto.ErrBadParam
	}
	if s.EnableValidation && !CheckOrderNum(order.Num) {
		s.log.Debug("OrderService: Save. Order num validation error")
		return dto.ErrBadOrderNum
	}

	exOrder, err := s.dbOrder.GetByNum(ctx, order.Num)
	if err == nil {
		if exOrder.UserID == order.UserID {
			return dto.ErrOrderRegistered
		} else {
			return dto.ErrOrderRegisteredByAnotherUser
		}
	} else if err != &model.NoRowFound {
		s.log.Error("OrderService: Save. Unexpected error", zap.Error(err))
		return err
	}
	modelOrder := s.mapOrderDTOtoModel(order)
	modelOrder.Status = "NEW"
	modelOrder.UploadAt = time.Now().Truncate(time.Microsecond)
	modelOrder.UpdatedAt = time.Now().Truncate(time.Microsecond)

	err = s.dbOrder.Save(ctx, modelOrder)
	if err != nil {
		s.log.Error("OrderService: Save. Can't save order",
			zap.Int("userID", order.UserID),
			zap.String("num", order.Num),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *OrderService) GetOrderList(ctx context.Context, userID int) ([]dto.Order, error) {
	if userID == 0 {
		s.log.Debug("OrderService: GetOrderList. got nil userID")
		return nil, dto.ErrBadParam
	}

	orderList, err := s.dbOrder.FindByUser(ctx, userID)
	if err != nil {
		s.log.Error("OrderService: GetOrderList. Can't get order list",
			zap.Int("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}

	resList := s.mapOrderListModelToDTO(orderList)
	return resList, nil
}
