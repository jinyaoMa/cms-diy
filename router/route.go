package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoutes(r *gin.Engine) {
	r.Use(Cors())

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", Login)
	}

	apiGroup := r.Group("/api")
	{
		api := apiGroup.Use(Auth())
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"test": "pong",
			})
		})
	}
}
