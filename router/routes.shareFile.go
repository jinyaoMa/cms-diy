package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ShareFileForm struct {
	Id       uint   `form:"id" binding:"required"`
	ExpireAt string `form:"expire" binding:"required"`
}

// @Summary ShareFile
// @Description Update the share state of a file
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param id formData uint true "ID"
// @Param expire formData string true "ExpireAt"
// @Success 200 {object} Json200Response "{"success":true,"data":{"file":{}}"
// @Failure 400 "ShareFileForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/shareFile [put]
func shareFile(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var form ShareFileForm
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

	expireAt, err := time.Parse(time.RFC3339Nano, form.ExpireAt)
	if err != nil {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid expire time format",
		})
		return
	}

	file.ShareExpiredAt = expireAt
	updateOk := model.SaveFile(&file)
	if !updateOk {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "fail to save file",
		})
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"file": file.APIFile,
			"share": JsonObject{
				"code":   file.ShareCode,
				"expire": file.ShareExpiredAt,
			},
		},
	})
}
