package model

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getFileType(fileInfo os.FileInfo) (fileType string) {
	if fileInfo.IsDir() {
		fileType = FILE_TYPE_DIRECTORY
	} else {
		fileType = FILE_TYPE_FILE
	}
	return
}

func autoFillDataWithAPath(file *File) error {
	file.Name = filepath.Base(file.APath)
	rpath, err := filepath.Rel(file.Workspace, file.APath)
	if err != nil {
		return err
	}
	if rpath != "." {
		file.Depth = strings.Count(rpath, string(filepath.Separator))
	}
	file.RPath = rpath
	file.IPath = parseSum256ToString(sha256.Sum256([]byte(file.APath)))
	file.Ext = filepath.Ext(file.APath)
	return nil
}

func refreshShareCode(file *File) {
	if file.Recycled == 1 {
		file.ShareCode = ""
		file.ShareExpiredAt = time.Now()
	} else if file.CanRefreshShareCode {
		if file.ShareExpiredAt.After(time.Now()) {
			file.ShareCode = generateShareCode(4)
		} else {
			file.ShareCode = ""
		}
	}
}
