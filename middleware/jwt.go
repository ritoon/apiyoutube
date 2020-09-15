package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"fmt"
	"net/http"
	"strings"
)

func VerifyJWT() gin.HandlerFunc {
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
		// Check JWT

		// sample token string taken from the New example
		hmacSampleSecret := []byte("my_secret_key")

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(res[1], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecret, nil
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
