package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name       string `gorm:"unique"`
	Permission string `gorm:"default:'*'"`
	Space      Size   `gorm:"default:0"`
	Code       string `gorm:"unique"`
	Valid      uint   `gorm:"check:valid IN (0, 1);default:1"`
	IsReserved uint   `gorm:"check:is_reserved IN (0, 1);default:0"` // default role to be reserved
}
type Roles []Role

func GetRoleByName(name string) (role Role, hasRole bool) {
	result := db.First(&role, "valid = 1 AND name = ?", name)
	hasRole = result.RowsAffected == 1
	return
}

func GetRoleById(id uint) (role Role, hasRole bool) {
	result := db.First(&role, "valid = 1 AND id = ?", id)
	hasRole = result.RowsAffected == 1
	return
}

func GetRoleByCode(code string) (role Role, hasRole bool) {
	result := db.First(&role, "valid = 1 AND code = ?", code)
	hasRole = result.RowsAffected == 1
	return
}
