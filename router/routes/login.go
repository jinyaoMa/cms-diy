package routes

import (
	"jinyaoma/cms-diy/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login() gin.HandlerFunc {
	return func (c *gin.Context) {
		var form LoginForm
		if bindFormPost(c, &form) != nil {
			return
		}
	
		user, hasUser := model.GetUserByAccountPassword(form.Account, form.Password)
		if !hasUser {
			c.JSON(http.StatusOK, gin.H{
				"error": "user password unmatched",
			})
			return
		}
	
		role, hasRole := model.GetRoleById(user.RoleID)
		if !hasRole {
			c.JSON(http.StatusOK, gin.H{
				"error": "user has no role assigned",
			})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"userid":     user.ID,
				"username":   user.Name,
				"role":       role.Name,
				"permission": role.Permission,
			},
		})
	}	
}
