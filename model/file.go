package model

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name           string
	APath          string `gorm:"index"` // absolute path
	RPath          string `gorm:"index"` // relative path (relative to user space folder)
	IPath          string `gorm:"index"` // id string for absolute path
	TPath          string `gorm:"index"` // thumbnail path
	Depth          int
	Type           string `gorm:"check:type IN ('directory', 'file')"`
	Ext            string
	Size           Size
	ShareCode      string
	ShareExpiredAt time.Time
	Recycled       uint `gorm:"check:recycled IN (0, 1);default:0"`
	UserID         uint
}
type Files []File
type APIFile struct {
	gorm.Model
	Name     string
	RPath    string
	IPath    string
	TPath    string
	Depth    int
	Type     string
	Ext      string
	Size     Size
	Recycled uint
}
type APIFiles []APIFile

const (
	FILE_TYPE_DIRECTORY string = "directory"
	FILE_TYPE_FILE      string = "file"
)

func FindFilesByUser(user User, depth int, offset int, limit int) (userFiles APIFiles, ok bool) {
	result := db.Model(&Files{}).
		Offset(offset).
		Limit(limit).
		Where("recycled = 0 AND depth = ? AND user_id = ?", depth, user.ID).
		Order("type, r_path").
		Find(&userFiles)
	ok = result.RowsAffected > 0
	return
}

func CreateUserSpaceFiles(userAccount string) (userFiles Files, err error) {
	err = NewUserSpace(userAccount, func(apath string, rpath string, depth int, fileInfo os.FileInfo) {
		var fileType string
		if fileInfo.IsDir() {
			fileType = FILE_TYPE_DIRECTORY
		} else {
			fileType = FILE_TYPE_FILE
		}
		userFiles = append(userFiles, File{
			Name:  filepath.Base(apath),
			IPath: fmt.Sprintf("%x", sha256.Sum256([]byte(apath))),
			APath: apath,
			RPath: rpath,
			Depth: depth,
			Type:  fileType,
			Ext:   filepath.Ext(apath),
			Size:  Size(fileInfo.Size()),
		})
	}, false)
	if err != nil {
		return userFiles, newError("User[" + userAccount + "] space initializing error")
	}
	return
}

func InitializeUserSpaceFiles(user User) (userFiles Files, err error) {
	err = NewUserSpace(user.Account, func(apath string, rpath string, depth int, fileInfo os.FileInfo) {
		var fileType string
		if fileInfo.IsDir() {
			fileType = FILE_TYPE_DIRECTORY
		} else {
			fileType = FILE_TYPE_FILE
		}
		userFiles = append(userFiles, File{
			Name:   filepath.Base(apath),
			IPath:  fmt.Sprintf("%x", sha256.Sum256([]byte(apath))),
			APath:  apath,
			RPath:  rpath,
			Depth:  depth,
			Type:   fileType,
			Ext:    filepath.Ext(apath),
			Size:   Size(fileInfo.Size()),
			UserID: user.ID,
		})
	}, true)
	if err != nil {
		return userFiles, newError("User[" + user.Account + "] space initializing error")
	}
	resultUserFiles := db.CreateInBatches(userFiles, 1000)
	if resultUserFiles.Error != nil {
		return userFiles, newError("User[" + user.Account + "] files init batches error")
	}
	return
}
