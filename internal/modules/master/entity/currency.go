package entity

type Currency struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

func (Currency) TableName() string { return "mst_currencies" }
