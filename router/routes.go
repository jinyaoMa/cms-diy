package router

import (
	"github.com/gin-gonic/gin"
)

type JsonObject map[string]interface{}
type Json200Response struct {
	Success bool       `json:"success"`
	Data    JsonObject `json:"data"`
}
type Json404Response struct {
	Error string `json:"error"`
}

func NewRoutes(r *gin.Engine) {
	r.Use(Cors())

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", Login)
	}

	apiGroup := r.Group("/api")
	{
		api := apiGroup.Use(Auth())
		api.GET("/test", Test)
	}
}
