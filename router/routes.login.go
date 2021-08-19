package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// @Summary Login
// @Description Login with account and password
// @Tags Before Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Param account formData string true "Account"
// @Param password formData string true "Password"
// @Success 200 {object} Json200Response "{"success":true,"data":{"userid":1,"username":"cms-diy","role":"admin","permission":"*","token":""}}"
// @Failure 400 "LoginForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var form LoginForm
	if bindFormPost(c, &form) != nil {
		return
	}

	user, hasUser := model.GetUserByAccountPassword(form.Account, form.Password)
	if !hasUser {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "user password unmatched",
		})
		return
	}

	role, hasRole := model.GetRoleById(user.RoleID)
	if !hasRole {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "user has no role assigned",
		})
		return
	}

	token, err := generateToken(c, user)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"userid":     user.ID,
			"username":   user.Name,
			"role":       role.Name,
			"permission": role.Permission,
			"token":      token,
		},
	})
}
