package entity

type Religion struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Religion) TableName() string { return "mst_religions" }
