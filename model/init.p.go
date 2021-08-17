package model

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

var (
	randSource rand.Source
)

func init_p() {
	randSource = rand.NewSource(time.Now().Unix())
}

func newError(s string) error {
	return errors.New(s)
}

func println(s string) {
	log.Println(s)
}

func generateToken(len int) string {
	r := rand.New(randSource)
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		var b int
		randLetterOrDigit := r.Intn(5)
		if randLetterOrDigit%2 == 0 { // letters
			randUpperOrLower := r.Intn(5)
			if randUpperOrLower%2 == 0 { // uppercase
				b = r.Intn(26) + 65
			} else { // lowercase
				b = r.Intn(26) + 97
			}
		} else { // digits
			b = r.Intn(10) + 48
		}
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func generateShareCode(len int) string {
	r := rand.New(randSource)
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		var b int
		randLetterOrDigit := r.Intn(5)
		if randLetterOrDigit%2 == 0 { // letters
			b = r.Intn(26) + 97
		} else { // digits
			b = r.Intn(10) + 48
		}
		bytes[i] = byte(b)
	}
	return string(bytes)
}
