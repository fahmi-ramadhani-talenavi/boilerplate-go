package entity

type SubMenu struct {
	ID       string `json:"id" gorm:"primaryKey"`
	ParentID string `json:"parent_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Order    int    `json:"order"`
}

func (SubMenu) TableName() string { return "sys_sub_menus" }
