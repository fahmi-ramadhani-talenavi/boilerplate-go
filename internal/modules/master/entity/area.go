package entity

type Area struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Area) TableName() string { return "mst_areas" }
