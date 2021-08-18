package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(Cors())

	apiGroup := router.Group("/api")
	{
		api := apiGroup.Use(Auth())
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"test": "pong",
			})
		})
		init_login(apiGroup)
	}
}

func Run() {
	router.Run(":" + SERVER_PORT)
}
