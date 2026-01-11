package entity

type EducationLevel struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func (EducationLevel) TableName() string { return "mst_education_levels" }
