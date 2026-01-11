package entity

type Branch struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (Branch) TableName() string { return "mst_branches" }
