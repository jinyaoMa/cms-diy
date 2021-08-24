package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupForm struct {
	Username *string `form:"username" binding:"required"`
	Account  *string `form:"account" binding:"required"`
	Password *string `form:"password" binding:"required"`
	Code     *string `form:"code" binding:"required"`
}

// @Summary Signup
// @Description Signup as a new user with role by invitation code
// @Tags Before Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param account formData string true "Account"
// @Param password formData string true "Password"
// @Param code formData string true "Invitation Code"
// @Success 200 {object} Json200Response "{"success":true,"data":{"userid":1,"username":"cms-diy","role":"admin","permission":"*","token":""}}"
// @Failure 400 "SignupForm binding error"
// @Failure 404 {object} Json404Response "{"error":"error msg"}"
// @Failure 500 "Token generating error"
// @Router /auth/signup [post]
func signup(c *gin.Context) {
	userCount, ok := model.GetActiveUsersCount()
	if !ok || userCount >= model.UserLimit {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "limited users for signup",
		})
		return
	}

	var form SignupForm
	if bindFormPost(c, &form) != nil {
		return
	}

	role, ok := model.GetRoleByCode(*form.Code)
	if !ok {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid invitation code",
		})
		return
	}

	userfiles, errCreateUserSpaceFiles := model.CreateUserSpaceFiles(*form.Account)
	if errCreateUserSpaceFiles != nil {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "invalid user account",
		})
		return
	}

	newUser := model.User{
		Name:     *form.Username,
		Account:  *form.Account,
		Password: *form.Password,
		RoleID:   role.ID,
		Files:    userfiles,
	}
	errCreateUser := model.CreateUser(&newUser)
	if errCreateUser != nil {
		c.JSON(http.StatusNotFound, Json404Response{
			Error: "fail to create user",
		})
		return
	}

	token, err := generateToken(c, newUser)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, Json200Response{
		Success: true,
		Data: JsonObject{
			"userid":     newUser.ID,
			"username":   newUser.Name,
			"role":       role.Name,
			"permission": role.Permission,
			"token":      token,
		},
	})
}
