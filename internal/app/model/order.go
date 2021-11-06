package model

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../service/mocks/mock_order_repository.go -package=mocks . OrderRepository
type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, orderID int) (*Order, error)
	GetByNum(ctx context.Context, num string) (*Order, error)
	UpdateStatus(ctx context.Context, userID int, num string, statusNew string) error
	FindByUser(ctx context.Context, userID int) ([]Order, error)
}

type Order struct {
	ID        int
	UserID    int
	Num       string
	Status    string
	UploadAt  time.Time
	UpdatedAt time.Time
}
