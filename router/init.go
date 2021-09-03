package router

import (
	"github.com/gin-gonic/gin"

	_ "jinyaoma/cms-diy/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title CMS_DIY (BACKEND)
// @version 0.0.1
// @description "try to be a NAS-like CMS"

// @contact.name Github Issues
// @contact.url https://github.com/jinyaoMa/cms-diy/issues

// @license.name MIT
// @license.url https://github.com/jinyaoMa/cms-diy/blob/main/LICENSE

// @securityDefinitions.apikey BearerIdAuth
// @in header
// @name Authorization

var (
	router *gin.Engine
)

func init() {
	if IS_RELEASE {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
}

func Run() {
	NewRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":" + SERVER_PORT)

	println("ROUTER RUN...")
}
