package entity

type TaxBracket struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	MinAmount  float64 `json:"min_amount"`
	MaxAmount  float64 `json:"max_amount"`
	Percentage float64 `json:"percentage"`
}

func (TaxBracket) TableName() string { return "mst_tax_brackets" }
