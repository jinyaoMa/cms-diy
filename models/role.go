package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name       string
	Permission string `gorm:"default:'CORE:1,SHARE:0,ADMIN:0'"`
	Space      uint64 `gorm:"default:0"`
}

type Roles []Role
