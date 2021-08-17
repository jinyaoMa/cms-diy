package model

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

func initRootUser() {
	hasRoot, _ := HasUserAccount(ROOT_ACCOUNT)
	hasRole, role := HasRoleName(ROOT_ROLENAME)
	if !hasRoot && hasRole {
		user := User{
			Name:     ROOT_USERNAME,
			Account:  ROOT_ACCOUNT,
			Password: ROOT_PASSWORD,
			JwtKey:   generateToken(JWT_KEY_SIZE),
			RoleID:   role.ID,
		}
		if CreateUser(&user) == nil {
			println("User root initialized")
		}

		var userFiles []File
		err := InitUserSpace(ROOT_ACCOUNT, func(apath string, fileInfo os.FileInfo) {
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
				Type:   fileType,
				Ext:    filepath.Ext(apath),
				Size:   Size(fileInfo.Size()),
				UserID: user.ID,
			})
		})
		if err != nil {
			println("User root space initializing error")
		}
		resultUserFiles := db.CreateInBatches(userFiles, 1000)
		if resultUserFiles.Error != nil {
			println("User files init batches error")
		}
	}
}
