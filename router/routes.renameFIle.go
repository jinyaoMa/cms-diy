package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type RenameFileForm struct {
	Id       uint   `form:"id" binding:"required"`
	Filename string `form:"filename" binding:"required"`
}

// @Summary RenameFile
// @Description Change name of a file
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id formData uint true "ID"
// @Param filename formData string true "Filename"
// @Success 200 {object} Json200Response "{"success":true,"data":{"file":{}}"
// @Failure 400 "RenameFileForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/renameFile [put]
func renameFile(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var form RenameFileForm
	if bindFormPost(c, &form) != nil {
		return
	}

	if !model.IsFileNameCharValid(form.Filename) {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid filename",
		})
		return
	}

	file, ok := model.GetFileByUserAndId(user, form.Id)
	if !ok {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid fileId",
		})
		return
	}

	currentDir := filepath.Dir(file.APath)
	newAPath := filepath.Join(currentDir, form.Filename)

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
