package model

import (
	"context"
	"time"
)

type OrderRepository interface {
	Save(ctx context.Context, order Order) error
	GetByID(ctx context.Context, orderID int) (*Order, error)
	GetByNum(ctx context.Context, num string) (*Order, error)
	UpdateStatus(ctx context.Context, user_id int, num string, statusNew string) error
	FindByUser(ctx context.Context, userID int) ([]Order, error)
}

type Order struct {
	Id         int
	User_id    int
	Num        string
	Status     string
	Upload_at  time.Time
	Updated_at time.Time
}
