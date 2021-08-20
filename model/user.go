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
	Active   uint `gorm:"check:active IN (0, 1);default:1"`
	RoleID   uint
	Role     Role
	Files    Files
}
type Users []User

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.JwtKey = generateToken(JWT_KEY_SIZE)
	return nil
}

func CreateUser(user *User) error {
	return db.Create(user).Error
}

func GetActiveUsersCount() (count int64, ok bool) {
	result := db.Model(&User{}).Where("active = 1").Count(&count)
	ok = result.Error == nil
	return
}

func GetUserById(id string) (user User, hasUser bool) {
	result := db.First(&user, "active = 1 AND id = ?", id)
	hasUser = result.RowsAffected == 1
	return
}

func GetUserByAccount(account string) (user User, hasUser bool) {
	result := db.First(&user, "active = 1 AND account = ?", account)
	hasUser = result.RowsAffected == 1
	return
}

func GetUserByAccountPassword(account string, password string) (user User, hasUser bool) {
	result := db.First(&user, "active = 1 AND account = ? AND password = ?", account, password)
	hasUser = result.RowsAffected == 1
	return
}
