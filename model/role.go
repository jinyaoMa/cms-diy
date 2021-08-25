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

func GetSpaceForUser(user User) (size Size, ok bool) {
	result := db.Raw("select space from roles where deleted_at IS NULL AND valid = 1 AND role_id = ?", user.RoleID).Scan(&size)
	ok = result.Error == nil
	return
}

func GetRoleByName(name string) (role Role, ok bool) {
	result := db.First(&role, "valid = 1 AND name = ?", name)
	ok = result.RowsAffected == 1
	return
}

func GetRoleById(id uint) (role Role, ok bool) {
	result := db.First(&role, "valid = 1 AND id = ?", id)
	ok = result.RowsAffected == 1
	return
}

func GetRoleByCode(code string) (role Role, ok bool) {
	result := db.First(&role, "valid = 1 AND code = ?", code)
	ok = result.RowsAffected == 1
	return
}
