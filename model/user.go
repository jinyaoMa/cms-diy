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

func CreateUser(user *User) error {
	return db.Create(user).Error
}

func GetUserById(id string) (user User, hasUser bool) {
	result := db.First(&user, id)
	hasUser = result.RowsAffected == 1
	return
}

func GetUserByAccount(account string) (user User, hasUser bool) {
	result := db.First(&user, "account = ?", account)
	hasUser = result.RowsAffected == 1
	return
}

func GetUserByAccountPassword(account string, password string) (user User, hasUser bool) {
	result := db.First(&user, "account = ? AND password = ?", account, password)
	hasUser = result.RowsAffected == 1
	return
}
