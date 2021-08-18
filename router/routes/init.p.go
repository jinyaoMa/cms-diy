package routes

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func newError(s string) error {
	return errors.New(s)
}

func println(s string) {
	log.Println(s)
}

func bindFormPost(c *gin.Context, form interface{}) error {
	err := c.ShouldBindWith(form, binding.FormPost)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return err
}