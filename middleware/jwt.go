package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"

	"fmt"
	"net/http"
	"strings"
)

type MyCustomClaims struct {
	Foo  string `json:"foo"`
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

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
		fmt.Println(res[1])
		fmt.Println([]byte(secret))
		fmt.Println([]byte("my_secret_key"))

		token, err := jwt.ParseWithClaims(res[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
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

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			fmt.Println(claims.Foo, claims.UUID)
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid user claims",
			})
		}
	}
}

func GenerateJWT(secret, uuid string) string {
	// mySigningKey := []byte("my_secret_key")

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		uuid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 1500,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	fmt.Printf("%v %v", ss, err)
	return ss
}
