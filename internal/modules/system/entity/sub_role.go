package entity

type SubRole struct {
	ID          string `json:"id" gorm:"primaryKey"`
	RoleID      string `json:"role_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (SubRole) TableName() string { return "sys_sub_roles" }
