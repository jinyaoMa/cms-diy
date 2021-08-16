package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Account  string
	Password string
	JwtKey   *string
	RoleID   uint
	Role     Role
	Files    Files
}

type Users []User
