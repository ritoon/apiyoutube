package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerfiyJWT() gin.HandleFunc {
	return func(ctx *gin.Context) {
		value := ctx.Request.Header("authorization")
		if len(value) == 0 {
			cxt.JSON(http.StatusUnauthorized, gin.H{
				"error": "need a JWT",
			})
			return
		}
		res := strings.Split(value, " ")
		if len(res) < 1 {
			cxt.JSON(http.StatusUnauthorized, gin.H{
				"error": "need a JWT",
			})
			return
		}
		// TODO verfiy JWT
	}
}
