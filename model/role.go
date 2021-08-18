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

func GetRoleByName(name string) (role Role, hasRole bool) {
	result := db.First(&role, "name = ?", name)
	hasRole = result.RowsAffected == 1
	return
}

func GetRoleById(id uint) (role Role, hasRole bool) {
	result := db.First(&role, id)
	hasRole = result.RowsAffected == 1
	return
}
