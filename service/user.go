package service

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ritoon/cours/discover/apiyoutube/db"
)

type ServiceUser struct {
	db db.DB
}

func NewUser(db db.DB) *ServiceUser {
	return &ServiceUser{
		db: db,
	}
}

func (su *ServiceUser) GetUser(ctx *gin.Context) {

	u, err := su.db.GetUser(ctx.Param("uuid"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, u)
}
