package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Test
// @Description Test Token
// @Tags After Authorization
// @accept plain
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} Json200Response "{"success":true,"data":{"user":"admin","isTokenValid":true}}"
// @Failure 500 "Token generating error"
// @Router /api/test [get]
func test(c *gin.Context) {
	user, claims := getUserClaimsFromAuth(c)

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"user":         user.Account,
			"isTokenValid": claims.ExpiresAt > time.Now().Unix(),
		},
	})
}
