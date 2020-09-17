package middleware

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"fmt"
	"net/http"
	"strings"
)

func VerifyJWT(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := ctx.Request.Header.Get("authorization")
		if len(value) == 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "need a JWT",
			})
			return
		}
		res := strings.Split(value, " ")
		if len(res) <= 1 || res[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "need a JWT",
			})
			return
		}

		_, err := parseJWT(secret, res[1])

		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "JWT not valid",
			})
			return
		}
	}
}

func parseJWT(secret, tokenValue string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("not able to cast JWT")
	}
	return claims, nil
}

type CustomClaims struct {
	UUID string
	jwt.StandardClaims
}

func GenerateJWT(secret, uuid string, exp time.Time) string {

	// Create the Claims
	claims := CustomClaims{
		uuid,
		jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	fmt.Printf("%v %v", ss, err)
	return ss
}
