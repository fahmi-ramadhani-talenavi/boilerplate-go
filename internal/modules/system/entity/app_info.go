package entity

import sharedentity "github.com/user/go-boilerplate/internal/shared/entity"

type AppInfo struct {
	sharedentity.Base
	Key   string `json:"key" gorm:"uniqueIndex"`
	Value string `json:"value"`
	Group string `json:"group"`
}

func (AppInfo) TableName() string { return "sys_settings" }
