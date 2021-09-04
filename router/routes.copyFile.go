package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type CopyFileForm struct {
	Id uint   `form:"id" binding:"required"`
	To string `form:"to" binding:"required"`
}

// @Summary CopyFile
// @Description Copy a file to destination path
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id formData uint true "File ID"
// @Param to formData string true "To"
// @Success 200 {object} Json200Response "{"success":true,"data":{"file":{}}"
// @Failure 400 "CopyFileForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/copyFile [post]
func copyFile(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var form CopyFileForm
	if bindFormPost(c, &form) != nil {
		return
	}

	file, ok := model.GetFileByUserAndId(user, form.Id)
	if !ok {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid fileId",
		})
		return
	}

	targetDir := filepath.Join(file.Workspace, form.To)
	if !model.IsPathCharValid(form.To) || !strings.HasPrefix(targetDir, file.Workspace) {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid path",
		})
		return
	}

	newAPath := filepath.Join(targetDir, file.Name)

	errRename := os.Rename(file.APath, newAPath)
	if errRename != nil {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "fail to rename apath",
		})
		return
	}

	file.APath = newAPath
	updateOk := model.SaveFile(&file)
	if !updateOk {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "fail to save filename",
		})
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"file": file.APIFile,
		},
	})
}
