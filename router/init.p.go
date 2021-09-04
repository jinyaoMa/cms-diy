package router

import (
	"errors"
	"fmt"
	"log"
)

func newError(s string) error {
	return errors.New(s)
}

func println(s string) {
	log.Println(s)
}

func parseNumberToString(num ...interface{}) string {
	return fmt.Sprintf("%d", num...)
}
