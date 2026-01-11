package entity

type Bank struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Bank) TableName() string { return "mst_banks" }
