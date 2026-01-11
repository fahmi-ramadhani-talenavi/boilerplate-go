package entity

type Citizenship struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Citizenship) TableName() string { return "mst_citizenships" }
