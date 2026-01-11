package entity

type MaritalStatus struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (MaritalStatus) TableName() string { return "mst_marital_statuses" }
