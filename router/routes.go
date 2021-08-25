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
		auth := authGroup.Use()

		auth.POST("/login", login)
		auth.POST("/signup", signup)
	}

	apiGroup := r.Group("/api")
	{
		api := apiGroup.Use(Auth())

		api.GET("/test", test)
		api.GET("/getNewToken", getNewToken)
		api.GET("/getFileList", getFileList)

		api.POST("/newFolder", newFolder)

		api.PUT("/renameFile", renameFile)
		api.PUT("/moveFile", moveFile)
		api.PUT("/shareFile", shareFile)
	}
}
