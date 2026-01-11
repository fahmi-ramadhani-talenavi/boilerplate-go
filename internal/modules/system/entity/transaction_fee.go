package entity

import sharedentity "github.com/user/go-boilerplate/internal/shared/entity"

type TransactionFee struct {
	sharedentity.Base
	TransactionType string  `json:"transaction_type"`
	FeeType         string  `json:"fee_type"`
	Amount          float64 `json:"amount"`
	Percentage      float64 `json:"percentage"`
}

func (TransactionFee) TableName() string { return "sys_transaction_fees" }
