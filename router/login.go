package router

import "github.com/gin-gonic/gin"

func init_login(rg *gin.RouterGroup) {
	rg.POST("/login", login)
}

func login(c *gin.Context) {

}
