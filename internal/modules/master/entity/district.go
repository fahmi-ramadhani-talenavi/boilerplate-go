package entity

type District struct {
	ID         string `json:"id" gorm:"primaryKey"`
	ProvinceID string `json:"province_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}

func (District) TableName() string { return "mst_districts" }
