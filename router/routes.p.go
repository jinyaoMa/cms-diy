package router

import (
	"fmt"
	"jinyaoma/cms-diy/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func bindFormPost(c *gin.Context, form interface{}) error {
	err := c.ShouldBindWith(form, binding.FormPost)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return err
}

func generateToken(c *gin.Context, user model.User) (token string, err error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Audience:  user.Name,
			ExpiresAt: now.AddDate(0, 1, 0).Unix(), // expire next month
			Id:        fmt.Sprintf("%d", user.ID),
			IssuedAt:  now.Unix(),
			Issuer:    model.ROOT_USERNAME,
			NotBefore: now.Unix(),
			Subject:   user.Name,
		},
	}

	jwToken := JWT{
		[]byte(user.JwtKey),
	}
	token, err = jwToken.CreateToken(claims)
	if err != nil {
		println(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	return
}
