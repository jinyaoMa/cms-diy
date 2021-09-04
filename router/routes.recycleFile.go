package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecycleFileForm struct {
	Id *uint `form:"id" binding:"required"` // 0 for at root of user workspace
}

// @Summary RecycleFile
// @Description Set a file/directory to be recycled
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id formData uint true "File/Directory ID (root - 0)"
// @Success 200 {object} Json200Response "{"success":true,"data":{"fileCount":0,"directoryCount":0}"
// @Failure 400 "RecycleFileForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/recycleFile [put]
func recycleFile(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var form RecycleFileForm
	if bindFormPost(c, &form) != nil {
		return
	}

	var fileCount, directoryCount uint
	var okRF bool
	if *form.Id > 0 {
		file, ok := model.GetFileByUserAndId(user, *form.Id)
		if !ok {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "invalid fileId",
			})
			return
		}

		fileCount, directoryCount, okRF = model.RecycleFile(file)
		if !okRF {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "fail to recycle file/directory",
			})
			return
		}
	} else {
		fileCount, directoryCount, okRF = model.RecycleFilesByUser(user)
		if !okRF {
			c.JSON(http.StatusNotFound, Json404Response{
				Error: "fail to recycle all user files/directories",
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
