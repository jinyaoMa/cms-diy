package main

import (
	"jinyaoma/cms-diy/model"
	"jinyaoma/cms-diy/router"
)

func main() {
	model.Run()
	router.Run()
}
