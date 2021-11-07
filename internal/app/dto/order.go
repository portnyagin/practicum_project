package dto

import "time"

type Order struct {
	Num      string    `json:"number"`
	UserID   int       `json:"-"`
	Status   string    `json:"status"`
	Accrual  int       `json:"accrual"`
	UploadAt time.Time `json:"upload_at"`
}
