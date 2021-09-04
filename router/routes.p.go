package router

import (
	"jinyaoma/cms-diy/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func bindQuery(c *gin.Context, query interface{}) (err error) {
	err = c.ShouldBind(query)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return
}

func bindFormPost(c *gin.Context, form interface{}) (err error) {
	err = c.ShouldBindWith(form, binding.FormPost)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	return
}

func getUserClaimsFromAuth(c *gin.Context) (user model.User, claims *JWTClaims) {
	var exists bool
	var temp interface{}

	temp, exists = c.Get("user")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user = temp.(model.User)

	temp, exists = c.Get("claims")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	claims = temp.(*JWTClaims)

	return
}

func generateToken(c *gin.Context, user model.User) (token string, err error) {
	now := time.Now().Unix()
	claims := JWTClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Audience:  user.Name,
			ExpiresAt: now + TOKEN_VALID_TIME_IN_SECOND,
			Id:        parseNumberToString(user.ID),
			IssuedAt:  now,
			Issuer:    model.ROOT_USERNAME,
			NotBefore: now,
			Subject:   user.Name,
		},
	}

	jwToken := JWT{
		[]byte(user.JwtKey),
	}
	token, err = jwToken.CreateToken(claims)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	return
}

func refreshToken(c *gin.Context, user model.User, claims JWTClaims) (token string, err error) {
	claims.NotBefore = time.Now().Unix()
	claims.ExpiresAt = claims.NotBefore + TOKEN_VALID_TIME_IN_SECOND

	jwToken := JWT{
		[]byte(user.JwtKey),
	}
	token, err = jwToken.CreateToken(claims)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	return
}
