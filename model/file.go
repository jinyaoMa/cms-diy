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
	APIFile
	APath          string `gorm:"index"` // absolute path
	Workspace      string `gorm:"index"` // user root path
	ShareCode      string
	ShareExpiredAt time.Time
	UserID         uint
}
type Files []File
type APIFile struct {
	gorm.Model
	Name     string
	RPath    string `gorm:"index"` // relative path (relative to user space folder)
	IPath    string `gorm:"index"` // id string for absolute path
	TPath    string `gorm:"index"` // thumbnail path
	Depth    int
	Type     string `gorm:"check:type IN ('directory', 'file')"`
	Ext      string
	Size     Size
	Recycled uint `gorm:"check:recycled IN (0, 1);default:0"`
}
type APIFiles []APIFile

const (
	FILE_TYPE_DIRECTORY string = "directory"
	FILE_TYPE_FILE      string = "file"
)

func (f *File) BeforeSave(tx *gorm.DB) (err error) {
	if f.ShareExpiredAt.After(time.Now()) {
		f.ShareCode = generateShareCode(4)
	} else {
		f.ShareCode = ""
	}
	return nil
}

func SaveFile(file *File) (ok bool) {
	result := db.Save(file)
	ok = result.RowsAffected == 1
	return
}

func GetFileByUserAndId(user User, id uint) (userFile File, ok bool) {
	result := db.Where("recycled = 0 AND user_id = ? AND id = ?", user.ID, id).First(&userFile)
	ok = result.RowsAffected == 1
	return
}

func FindAPIFilesByUser(user User, depth int, offset int, limit int) (userFiles APIFiles, ok bool) {
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
	err = NewUserSpace(userAccount, func(apath string, rpath string, depth int, workspace string, fileInfo os.FileInfo) {
		var fileType string
		if fileInfo.IsDir() {
			fileType = FILE_TYPE_DIRECTORY
		} else {
			fileType = FILE_TYPE_FILE
		}
		userFiles = append(userFiles, File{
			APIFile: APIFile{
				Name:  filepath.Base(apath),
				RPath: rpath,
				IPath: fmt.Sprintf("%x", sha256.Sum256([]byte(apath))),
				Depth: depth,
				Type:  fileType,
				Ext:   filepath.Ext(apath),
				Size:  Size(fileInfo.Size()),
			},
			APath:     apath,
			Workspace: workspace,
		})
	}, false)
	if err != nil {
		return userFiles, newError("User[" + userAccount + "] space initializing error")
	}
	return
}

func InitializeUserSpaceFiles(user User) (userFiles Files, err error) {
	err = NewUserSpace(user.Account, func(apath string, rpath string, depth int, workspace string, fileInfo os.FileInfo) {
		var fileType string
		if fileInfo.IsDir() {
			fileType = FILE_TYPE_DIRECTORY
		} else {
			fileType = FILE_TYPE_FILE
		}
		userFiles = append(userFiles, File{
			APIFile: APIFile{
				Name:  filepath.Base(apath),
				RPath: rpath,
				IPath: fmt.Sprintf("%x", sha256.Sum256([]byte(apath))),
				Depth: depth,
				Type:  fileType,
				Ext:   filepath.Ext(apath),
				Size:  Size(fileInfo.Size()),
			},
			APath:     apath,
			Workspace: workspace,
			UserID:    user.ID,
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
