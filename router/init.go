package router

import (
	"jinyaoma/cms-diy/router/routes"
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

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", routes.Login())
	}

	apiGroup := router.Group("/api")
	{
		api := apiGroup.Use(Auth())
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"test": "pong",
			})
		})
	}
}

func Run() {
	router.Run(":" + SERVER_PORT)
}
