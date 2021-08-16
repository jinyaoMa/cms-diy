package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name           string
	Path           string `gorm:"not null;index"`
	Type           string `gorm:"check:type IN ('directory', 'file');not null"`
	Extension      string
	Size           uint64
	Recycled       uint `gorm:"check:recycled IN (0, 1);default:0"`
	ShareCode      string
	ShareExpiredAt time.Time
	UserID         uint
}

type Files []File
