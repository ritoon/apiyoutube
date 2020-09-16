package service

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"apiyoutube/db"
	"apiyoutube/model"
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

func (su *ServiceUser) GetListUser(ctx *gin.Context) {
	us, err := su.db.GetListUser()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}
	ctx.JSON(http.StatusOK, us)
}

func (su *ServiceUser) CreateUser(ctx *gin.Context) {
	var u model.User
	err := ctx.BindJSON(&u)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parsing user"})
		return
	}
	passHash := sha256.Sum256([]byte(u.Pass))
	u.Pass = base64.StdEncoding.EncodeToString(passHash[:])
	err = su.db.AddUser(&u)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (su *ServiceUser) UpdateUser(ctx *gin.Context) {
	var u model.User
	err := ctx.BindJSON(&u)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parsing user"})
		return
	}
	passHash := sha256.Sum256([]byte(u.Pass))
	u.Pass = string(passHash[:])
	err = su.db.UpdateUser(ctx.Param("uuid"), u)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (su *ServiceUser) DeleteUser(ctx *gin.Context) {
	err := su.db.DeleteUser(ctx.Param("uuid"))
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
