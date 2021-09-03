package router

import (
	"fmt"
	"jinyaoma/cms-diy/model"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	SigningKey []byte
}
type JWTClaims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

var (
	TokenExpired     error = newError("Token is expired")
	TokenNotValidYet error = newError("Token not active yet")
	TokenMalformed   error = newError("That's not even a token")
	TokenInvalid     error = newError("Couldn't handle this token:")
)

func (j *JWT) CreateToken(claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func Index(r *gin.Engine) gin.HandlerFunc {
	indexRegexp, err := regexp.Compile(`.*/$`)
	if err != nil {
		panic("Failed to compile regexp in middleware Index")
	}

	return func(c *gin.Context) {
		str := indexRegexp.FindString(c.Request.URL.Path)
		if str == "" {
			c.Next()
		} else {
			index := c.Request.URL.Path + "index.html"
			println(index)
			c.Redirect(http.StatusMovedPermanently, index)
			return
		}
	}
}

func Auth() gin.HandlerFunc {
	bearerRegexp, err := regexp.Compile(`^Bearer (\d+) (.+)$`)
	if err != nil {
		panic("Failed to compile regexp in middleware Auth")
	}

	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		matches := bearerRegexp.FindStringSubmatch(authorization)
		if matches == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userId := matches[1]
		user, ok := model.GetUserById(userId)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := matches[2]
		jwToken := JWT{
			SigningKey: []byte(user.JwtKey),
		}
		claims, err := jwToken.ParseToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if userId != fmt.Sprintf("%d", claims.UserID) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.ExpiresAt < time.Now().Unix() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
		c.Set("claims", claims)
		c.Next()
	}
}
