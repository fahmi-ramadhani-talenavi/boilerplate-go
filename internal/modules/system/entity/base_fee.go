package entity

import sharedentity "github.com/user/go-boilerplate/internal/shared/entity"

type BaseFee struct {
	sharedentity.Base
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

func (BaseFee) TableName() string { return "sys_base_fees" }
