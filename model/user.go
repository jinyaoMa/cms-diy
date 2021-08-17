package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Account  string `gorm:"unique"`
	Password string
	JwtKey   string
	RoleID   uint
	Role     Role
	Files    Files
}
type Users []User

func HasUserAccount(account string) (bool, User) {
	var user User
	result := db.First(&user, "account = ?", account)
	return result.RowsAffected == 1, user
}

func CreateUser(user *User) error {
	return db.Create(user).Error
}