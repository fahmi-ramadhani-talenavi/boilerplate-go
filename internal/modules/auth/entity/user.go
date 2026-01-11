package entity

import (
	sharedentity "github.com/user/go-boilerplate/internal/shared/entity"
)

// User represents a system user for authentication.
type User struct {
	sharedentity.Base
	Email          string  `json:"email" gorm:"uniqueIndex;not null"`
	Password       string  `json:"-" gorm:"not null"`
	FullName       string  `json:"full_name" gorm:"not null"`
	EmployeeNumber *string `json:"employee_number,omitempty"`
	BranchID       *string `json:"branch_id,omitempty" gorm:"type:uuid"`
	DivisionID     *string `json:"division_id,omitempty" gorm:"type:uuid"`
	DivisionName   *string `json:"division_name,omitempty"`
	RoleID         *string `json:"role_id,omitempty" gorm:"type:uuid"`
	SubRoleID      *string `json:"sub_role_id,omitempty" gorm:"type:uuid"`
	ProfileID      *string `json:"profile_id,omitempty"`
	CompanyProfID  *string `json:"company_profile_id,omitempty" gorm:"column:company_profile_id"`
	IsActive       bool    `json:"is_active" gorm:"default:true"`
}

// TableName returns the database table name.
func (User) TableName() string {
	return "sys_users"
}
