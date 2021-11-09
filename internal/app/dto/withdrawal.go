package dto

import "time"

type Withdrawal struct {
	OrderNum    string    `json:"order"`
	Amount      float32   `json:"sum"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}

type Withdraw struct {
	OrderNum string  `json:"order"`
	Amount   float32 `json:"sum"`
}

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
