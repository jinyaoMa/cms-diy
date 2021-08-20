package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExtendTokenExpireTimeForm struct {
	Username string `form:"username" binding:"required"`
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
	Code     string `form:"code" binding:"required"`
}

// @Summary ExtendTokenExpireTime
// @Description Make up a new token to extend expire time
// @Tags After Authorization
// @accept plain
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @Success 200 {object} Json200Response "{"success":true,"data":{"token":""}}"
// @Failure 500 "Token generating error"
// @Router /api/extendTokenExpireTime [get]
func extendTokenExpireTime(c *gin.Context) {
	user, claims := getUserClaimsFromAuth(c)
	token, err := refreshToken(c, user, *claims)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"token": token,
		},
	})
}
