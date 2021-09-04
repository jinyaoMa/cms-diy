package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type NewFolderForm struct {
	Id      *uint  `form:"id" binding:"required"` // 0 for at root of user workspace
	DirName string `form:"name" binding:"required"`
}

// @Summary NewFolder
// @Description Change name of a file
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id formData uint true "Destination ID (root - 0)"
// @Param name formData string true "Directory Name"
// @Success 200 {object} Json200Response "{"success":true,"data":{"file":{}}"
// @Failure 400 "NewFolderForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/newFolder [post]
func newFolder(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var form NewFolderForm
	if bindFormPost(c, &form) != nil {
		return
	}

	if model.IsDirectoryNamedCache(form.DirName) || model.IsDirectoryNamedRecycle(form.DirName) {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "dirname '" + form.DirName + "' is reserved",
		})
		return
	}

	if !model.IsFileNameCharValid(form.DirName) {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid folder path",
		})
		return
	}

	var des model.File
	var okDes bool
	if *form.Id > 0 {
		des, okDes = model.GetFileByUserAndId(user, *form.Id)
		if !okDes {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "invalid fileId",
			})
			return
		}
		if des.Type != model.FILE_TYPE_DIRECTORY {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "destination not a folder",
			})
			return
		}
	}

	workspace, ok := model.GetValidUserWorkspace(user.Account, 0) // folder size 0
	if !ok {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "limited workspace",
		})
		return
	}

	newAPath := filepath.Join(workspace, des.RPath, form.DirName)

	err := os.Mkdir(newAPath, os.ModeDir)
	if os.IsExist(err) {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "folder exists",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "fail to make directory",
		})
		return
	}

	directory := model.File{
		APIFile: model.APIFile{
			Type: model.FILE_TYPE_DIRECTORY,
		},
		APath:     newAPath,
		Workspace: workspace,
		UserID:    user.ID,
	}
	okCreateFile := model.CreateFile(&directory)
	if !okCreateFile {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "fail to create directory",
		})
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"file": directory.APIFile,
		},
	})
}
