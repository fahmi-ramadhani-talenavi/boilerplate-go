package entity

type TaxGroup struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (TaxGroup) TableName() string { return "mst_tax_groups" }
