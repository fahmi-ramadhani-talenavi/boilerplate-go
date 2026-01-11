package entity

type Gender struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Gender) TableName() string { return "mst_genders" }
