package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name           string
	IPath          string `gorm:"unique"`
	APath          string `gorm:"index"`
	Type           string `gorm:"check:type IN ('directory', 'file')"`
	Ext            string `gorm:"check:ext LIKE '.%'"`
	Size           Size
	ShareCode      string
	ShareExpiredAt time.Time
	Recycled       uint `gorm:"check:recycled IN (0, 1);default:0"`
	UserID         uint
}
type Files []File

const (
	FILE_TYPE_DIRECTORY string = "directory"
	FILE_TYPE_FILE      string = "file"
)
