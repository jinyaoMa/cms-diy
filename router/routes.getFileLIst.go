package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetFileListQuery struct {
	Depth  *int `form:"depth" binding:"required"`
	Offset *int `form:"offset" binding:"required"`
	Limit  *int `form:"limit" binding:"required"`
}

// @Summary GetFileList
// @Description Get file list of user space
// @Tags After Authorization
// @accept plain
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Param depth query int true "Depth"
// @Param offset query int true "Offset"
// @Param limit query int true "Limit"
// @Success 200 {object} Json200Response "{"success":true,"data":{"files":[]}"
// @Failure 400 "GetFileListQuery binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /api/getFileList [get]
func getFileList(c *gin.Context) {
	user, _ := getUserClaimsFromAuth(c)

	var query GetFileListQuery
	if bindQuery(c, &query) != nil {
		return
	}

	files, ok := model.FindAPIFilesByUser(user, *query.Depth, *query.Offset, *query.Limit)
	if !ok {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "no files in user space",
		})
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"files": files,
		},
	})
}
