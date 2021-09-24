package model

import (
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

type File struct {
	APIFile
	APath               string `gorm:"index"` // absolute path
	Workspace           string `gorm:"index"` // user root path
	ShareCode           string
	ShareExpiredAt      time.Time
	CanRefreshShareCode bool `gorm:"-"`
	UserID              uint
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

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	autoFillDataWithAPath(f)
	refreshShareCode(f)
	return nil
}

func (f *File) BeforeUpdate(tx *gorm.DB) (err error) {
	autoFillDataWithAPath(f)
	refreshShareCode(f)
	return nil
}

func (f *File) BeforeSave(tx *gorm.DB) (err error) {
	autoFillDataWithAPath(f)
	refreshShareCode(f)
	return nil
}

func DeleteFilesByUser(user User) (fileCount uint, directoryCount uint, ok bool) {
	var files Files
	resultFind := db.
		Where("recycled = 1 AND user_id = ?", user.ID).
		Order("type desc, depth desc").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Remove(GenerateRecycleAPath(f))
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultDelete := db.Where("id IN ?", ids).Delete(Files{})
	ok = resultDelete.RowsAffected == resultFind.RowsAffected
	return
}

func RestoreFilesByUser(user User) (fileCount uint, directoryCount uint, ok bool) {
	var files Files
	resultFind := db.
		Where("recycled = 1 AND user_id = ?", user.ID).
		Order("type, depth").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Rename(GenerateRecycleAPath(f), f.APath)
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultUpdate := db.Table("files").Where("id IN ?", ids).Update("recycled", 0)
	ok = resultUpdate.RowsAffected == resultFind.RowsAffected
	return
}

func RecycleFilesByUser(user User) (fileCount uint, directoryCount uint, ok bool) {
	var files Files
	resultFind := db.
		Where("recycled = 0 AND user_id = ?", user.ID).
		Order("type desc, depth desc").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Rename(f.APath, GenerateRecycleAPath(f))
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultUpdate := db.Table("files").Where("id IN ?", ids).Update("recycled", 1)
	ok = resultUpdate.RowsAffected == resultFind.RowsAffected
	return
}

func DeleteFile(file File) (fileCount uint, directoryCount uint, ok bool) {
	apath := file.APath
	if os.PathSeparator == '\\' {
		apath = strings.ReplaceAll(file.APath, "\\", "\\\\")
	}

	var files Files
	resultFind := db.
		Where("recycled = 1 AND a_path LIKE ? AND user_id = ?", apath+"%%", file.UserID).
		Order("type desc, depth desc").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Remove(GenerateRecycleAPath(f))
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultDelete := db.Where("id IN ?", ids).Delete(Files{})
	ok = resultDelete.RowsAffected == resultFind.RowsAffected
	return
}

func RestoreFile(file File) (fileCount uint, directoryCount uint, ok bool) {
	apath := file.APath
	if os.PathSeparator == '\\' {
		apath = strings.ReplaceAll(file.APath, "\\", "\\\\")
	}

	var files Files
	resultFind := db.
		Where("recycled = 1 AND a_path LIKE ? AND user_id = ?", apath+"%%", file.UserID).
		Order("type, depth").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Rename(GenerateRecycleAPath(f), f.APath)
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultUpdate := db.Table("files").Where("id IN ?", ids).Update("recycled", 0)
	ok = resultUpdate.RowsAffected == resultFind.RowsAffected
	return
}

func RecycleFile(file File) (fileCount uint, directoryCount uint, ok bool) {
	apath := file.APath
	if os.PathSeparator == '\\' {
		apath = strings.ReplaceAll(file.APath, "\\", "\\\\")
	}

	var files Files
	resultFind := db.
		Where("recycled = 0 AND a_path LIKE ? AND user_id = ?", apath+"%%", file.UserID).
		Order("type desc, depth desc").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Rename(f.APath, GenerateRecycleAPath(f))
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultUpdate := db.Table("files").Where("id IN ?", ids).Update("recycled", 1)
	ok = resultUpdate.RowsAffected == resultFind.RowsAffected
	return
}

func MoveFileTo(file File, to File) (fileCount uint, directoryCount uint, ok bool) {
	apath := file.APath
	if os.PathSeparator == '\\' {
		apath = strings.ReplaceAll(file.APath, "\\", "\\\\")
	}

	var files Files
	resultFind := db.
		Where("recycled = 0 AND a_path LIKE ? AND user_id = ?", apath+"%%", file.UserID).
		Order("type desc, depth desc").
		Find(&files)
	if resultFind.RowsAffected < 1 {
		ok = false
		return
	}

	ids := []uint{}
	for _, f := range files {
		err := os.Rename(f.APath, GenerateRecycleAPath(f))
		if err != nil {
			ok = false
			return
		}
		if f.Type == FILE_TYPE_FILE {
			fileCount++
		} else if f.Type == FILE_TYPE_DIRECTORY {
			directoryCount++
		}
		ids = append(ids, f.ID)
	}
	resultUpdate := db.Table("files").Where("id IN ?", ids).Update("recycled", 1)
	ok = resultUpdate.RowsAffected == resultFind.RowsAffected
	return
}

func CreateFile(file *File) (ok bool) {
	result := db.Create(file)
	ok = result.RowsAffected == 1
	return
}

func GetUsedSpaceForUser(user User) (size Size, ok bool) {
	result := db.Raw("select sum(size) from files where deleted_at IS NULL AND user_id = ?", user.ID).Scan(&size)
	ok = result.Error == nil
	return
}

func SaveFile(file *File) (ok bool) {
	result := db.Save(file)
	ok = result.RowsAffected == 1
	return
}

func GetRecycledFileByUserAndId(user User, id uint) (userFile File, ok bool) {
	result := db.
		Where("recycled = 1 AND user_id = ? AND id = ?", user.ID, id).
		First(&userFile)
	ok = result.RowsAffected == 1
	return
}

func GetFileByUserAndId(user User, id uint) (userFile File, ok bool) {
	result := db.
		Where("recycled = 0 AND user_id = ? AND id = ?", user.ID, id).
		First(&userFile)
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
	if IsDirectoryNamedCache(userAccount) || IsDirectoryNamedRecycle(userAccount) {
		return nil, newError("User[" + userAccount + "] cannot be named '" + DIRECTORY_CACHE_NAME + "' or '" + DIRECTORY_RECYCLE_NAME + "'")
	}
	err = NewUserSpace(userAccount, func(apath string, rpath string, depth int, workspace string, fileInfo os.FileInfo) {
		fileType := getFileType(fileInfo)
		userFiles = append(userFiles, File{
			APIFile: APIFile{
				Type: fileType,
				Size: Size(fileInfo.Size()),
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
		fileType := getFileType(fileInfo)
		userFiles = append(userFiles, File{
			APIFile: APIFile{
				Type: fileType,
				Size: Size(fileInfo.Size()),
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
