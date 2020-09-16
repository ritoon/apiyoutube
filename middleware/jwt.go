package middleware

import (
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

		token, err := jwt.Parse(res[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Invalid user token: %s, err: %s", res[1], err.Error()),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["foo"], claims["nbf"])
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid user claims",
			})
		}
	}
}
