package model

import "time"

type Operation struct {
	ID            int
	AccountID     int
	OrderID       int
	OrderNum      string
	OperationType string
	Amount        float32
	ProcessedAt   time.Time
}

const OperationDebit = "DEBIT"
const OperationCredit = "CREDIT"
