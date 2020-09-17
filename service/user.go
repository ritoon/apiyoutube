package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"apiyoutube/db"
	"apiyoutube/middleware"
	"apiyoutube/model"
	"apiyoutube/util"
)

type ServiceUser struct {
	db        db.DB
	SecretJWT string
}

func NewUser(db db.DB, secret string) *ServiceUser {
	return &ServiceUser{
		db:        db,
		SecretJWT: secret,
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
	u.Pass = util.Hash(u.Pass)
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
	u.Pass = util.Hash(u.Pass)
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

func (su *ServiceUser) LoginUser(ctx *gin.Context) {
	var l model.Login
	err := ctx.BindJSON(&l)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "parsing user"})
		return
	}

	u, err := su.db.GetUserByEmail(l.Email)
	if err != nil {
		log.Printf("service: try to get into the db %v", err)
		ctx.JSON(err.(db.ErrorDB).Code, gin.H{"error": "db"})
		return
	}

	if !util.HashValid(l.Pass, u.Pass) {
		log.Printf("service: the hash is not valid")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user/pass"})
		return
	}

	jwtValue := middleware.GenerateJWT(su.SecretJWT, u.UUID, time.Now().Add(time.Hour*3))

	ctx.JSON(http.StatusOK, gin.H{"jwt": jwtValue})
}
