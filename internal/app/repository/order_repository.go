package repository

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/database/query"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
	"go.uber.org/zap"
)

type OrderRepositoryImpl struct {
	h basedbhandler.DBHandler
	l *infrastructure.Logger
}

func NewOrderRepository(dbHandler basedbhandler.DBHandler, log *infrastructure.Logger) (model.OrderRepository, error) {
	var target OrderRepositoryImpl
	if dbHandler == nil {
		return nil, errors.New("can't init order repository")
	}
	target.h = dbHandler
	target.l = log
	return &target, nil
}

func (or *OrderRepositoryImpl) Save(ctx context.Context, order model.Order) error {
	err := or.h.Execute(ctx, query.CreateOrder, order.User_id, order.Num, order.Status, order.Upload_at, order.Updated_at)
	return err
}

func (or *OrderRepositoryImpl) GetByID(ctx context.Context, orderID int) (*model.Order, error) {
	var res model.Order
	row, err := or.h.QueryRow(ctx, query.GetOrderByID, orderID)
	if err != nil {
		or.l.Error("request error", zap.String("query", query.GetOrderByID), zap.Int("user_id", orderID), zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.Id, &res.User_id, &res.Num, &res.Status, &res.Upload_at, &res.Updated_at)
	if err != nil {
		or.l.Error("scan rows error", zap.String("query", query.GetOrderByID), zap.Int("user_id", orderID), zap.Error(err))
	}
	return &res, nil
}

func (or *OrderRepositoryImpl) GetByNum(ctx context.Context, num string) (*model.Order, error) {
	var res model.Order
	row, err := or.h.QueryRow(ctx, query.GetOrderByNum, num)
	if err != nil {
		or.l.Error("request error", zap.String("query", query.GetOrderByNum), zap.String("user_id", num), zap.Error(err))
		return nil, err
	}
	err = row.Scan(&res.Id, &res.User_id, &res.Num, &res.Status, &res.Upload_at, &res.Updated_at)
	if err != nil {
		or.l.Error("scan rows error", zap.String("query", query.GetOrderByNum), zap.String("user_id", num), zap.Error(err))
	}
	return &res, nil
}

func (or *OrderRepositoryImpl) UpdateStatus(ctx context.Context, user_id int, num string, statusNew string) error {
	err := or.h.Execute(ctx, query.UpdateOrderStatus, user_id, num, statusNew)
	return err
}

// Проверить запрос и соответсвие полей
func (or *OrderRepositoryImpl) FindByUser(ctx context.Context, userID int) ([]model.Order, error) {
	rows, err := or.h.Query(ctx, query.FindOrdersByUser, userID)
	var resArray []model.Order

	if err != nil {
		or.l.Error("request error", zap.String("query", query.FindOrdersByUser), zap.Int("user_id", userID), zap.Error(err))
	}

	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.Id, &o.Num, &o.User_id, &o.Status, &o.Upload_at, &o.Updated_at)
		if err != nil {
			or.l.Error("scan rows error", zap.String("query", query.FindOrdersByUser), zap.Int("user_id", userID), zap.Error(err))
			break
		}
		resArray = append(resArray, o)
	}
	return resArray, nil
}
