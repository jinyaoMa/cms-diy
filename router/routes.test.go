package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Login
// @Description Login with account and password
// @Tags After Authorization
// @accept application/x-www-form-urlencoded
// @Produce  json
// @Security BearerIdAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} string "Pass"
// @Router /api/test [get]
func Test(c *gin.Context) {
	c.JSON(http.StatusOK, "pass")
}
