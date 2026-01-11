package entity

import sharedentity "github.com/user/go-boilerplate/internal/shared/entity"

type BankFee struct {
	sharedentity.Base
	BankID  string  `json:"bank_id"`
	FeeType string  `json:"fee_type"`
	Amount  float64 `json:"amount"`
}

func (BankFee) TableName() string { return "sys_bank_fees" }
