package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Permission string `gorm:"default:'*'"`
	Space      Size   `gorm:"default:0"`
	Code       string `gorm:"unique"`
}
type Roles []Role

func HasRoleName(name string) (bool, Role) {
	var role Role
	result := db.First(&role, "name = ?", name)
	return result.RowsAffected == 1, role
}
