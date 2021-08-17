package router

import (
	"jinyaoma/cms-diy/models"

	"github.com/gin-gonic/gin"
)

func init() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(Cors())

	api := router.Group("/api").Use(Auth())
	{
		
	}

	router.Run(SERVER_PORT)
}