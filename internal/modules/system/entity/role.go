package entity

type Role struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (Role) TableName() string { return "sys_roles" }
