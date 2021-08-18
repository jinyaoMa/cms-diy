package router

import (
	"errors"
	"log"
)

func newError(s string) error {
	return errors.New(s)
}

func println(s string) {
	log.Println(s)
}
