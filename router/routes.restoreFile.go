package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestoreFileForm struct {
	Id *uint `form:"id" binding:"required"` // 0 for at root of user workspace
}

// @Summary RestoreFile
// @Description Restore a recycled file/directory
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id formData uint true "File/Directory ID (root - 0)"
// @Success 200 {object} Json200Response "{"success":true,"data":{"fileCount":0,"directoryCount":0}"
// @Failure 400 "RestoreFileForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/restoreFile [put]
func restoreFile(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var form RestoreFileForm
	if bindFormPost(c, &form) != nil {
		return
	}

	var fileCount, directoryCount uint
	var okRF bool
	if *form.Id > 0 {
		file, ok := model.GetRecycledFileByUserAndId(user, *form.Id)
		if !ok {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "invalid fileId",
			})
			return
		}

		fileCount, directoryCount, okRF = model.RestoreFile(file)
		if !okRF {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "fail to restore file/directory",
			})
			return
		}
	} else {
		fileCount, directoryCount, okRF = model.RestoreFilesByUser(user)
		if !okRF {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "fail to restore all user files/directories",
			})
			return
		}
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"fileCount":      fileCount,
			"directoryCount": directoryCount,
		},
	})
}
