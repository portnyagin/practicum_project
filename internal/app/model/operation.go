package model

import "time"

type Operation struct {
	ID            int
	AccountID     int
	OrderID       int
	OrderNum      string
	OperationType string
	Amount        int
	ProcessedAt   time.Time
}

const OperationDebit = "DEBIT"
const OperationCredit = "CREDIT"
