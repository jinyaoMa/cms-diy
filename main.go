package main

import (
	"jinyaoma/cms-diy/model"
	"jinyaoma/cms-diy/router"
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

func main() {
	model.Run()
	router.Run()
}
