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
	UpdateStatus(ctx context.Context, order *Order) error
	FindByUser(ctx context.Context, userID int) ([]Order, error)
	LockOrder(ctx context.Context, OrderNum string) (*Order, error)
	FindNotProcessed(ctx context.Context) ([]Order, error)
}

type Order struct {
	ID        int
	UserID    int
	Num       string
	Status    string
	UploadAt  time.Time
	UpdatedAt time.Time
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusRegistered = "REGISTERED"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
)

func IsFinal(status string) bool {
	if status == OrderStatusInvalid || status == OrderStatusProcessed {
		return true
	}
	return false
}
