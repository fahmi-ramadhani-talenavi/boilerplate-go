package entity

type Province struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Province) TableName() string { return "mst_provinces" }
