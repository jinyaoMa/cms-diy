package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name           string
	APath          string `gorm:"index"` // absolute path
	IPath          string `gorm:"index"` // id string for absolute path
	TPath          string `gorm:"index"` // thumbnail path
	Type           string `gorm:"check:type IN ('directory', 'file')"`
	Ext            string
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
