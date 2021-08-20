package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Test
// @Description Test Token
// @Tags After Authorization
// @accept plain
// @Produce  plain
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} string "Pass"
// @Router /api/test [get]
func Test(c *gin.Context) {
	c.JSON(http.StatusOK, "pass")
}
