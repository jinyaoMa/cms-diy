package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteFileQuery struct {
	Id *uint `form:"id" binding:"required"` // 0 for at root of user workspace
}

// @Summary DeleteFile
// @Description Delete a file/directory permanently
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id query uint true "File/Directory ID (root - 0)"
// @Success 200 {object} Json200Response "{"success":true,"data":{"files":[]}"
// @Failure 400 "DeleteFileForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/deleteFile [delete]
func deleteFile(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var query DeleteFileQuery
	if bindQuery(c, &query) != nil {
		return
	}

	var deletedFiles model.APIFiles
	var okRF bool
	if *query.Id > 0 {
		file, ok := model.GetRecycledFileByUserAndId(user, *query.Id)
		if !ok {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "invalid fileId",
			})
			return
		}

		deletedFiles, okRF = model.DeleteFile(file)
		if !okRF {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "fail to delete file/directory",
			})
			return
		}
	} else {
		deletedFiles, okRF = model.DeleteFilesByUser(user)
		if !okRF {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "fail to delete all user files/directories",
			})
			return
		}
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"files": deletedFiles,
		},
	})
}
