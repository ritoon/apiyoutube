package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	kafka "github.com/segmentio/kafka-go"

	"apiyoutube/db"
	"apiyoutube/middleware"
	"apiyoutube/model"
	"apiyoutube/queue"
	"apiyoutube/util"
)

type ServiceUser struct {
	db          db.DB
	cache       *redis.Client
	secretJWT   string
	kafkaWriter *kafka.Writer
}

func NewUser(db db.DB, cache *redis.Client, writer *kafka.Writer, secret string) *ServiceUser {
	return &ServiceUser{
		db:          db,
		secretJWT:   secret,
		cache:       cache,
		kafkaWriter: writer,
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
		errval, ok := err.(db.ErrorDB)
		if ok {
			log.Println(errval)
			ctx.JSON(errval.Code, gin.H{"error": "db"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "db"})
		return
	}

	if !util.HashValid(l.Pass, u.Pass) {
		log.Printf("service: the hash is not valid")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user/pass"})
		return
	}
	msg := queue.MessageLogin{
		UUID:        u.UUID,
		UserAgent:   ctx.Request.UserAgent(),
		CallContext: "login",
		IP:          ctx.ClientIP(),
		EventTime:   time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("service: parse error ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user/pass"})
		return
	}

	su.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(u.UUID),
			Value: []byte(data),
		})

	// create JWT
	jwtValue := middleware.GenerateJWT(su.secretJWT, u.UUID, time.Now().Add(time.Hour*3))

	ctx.JSON(http.StatusOK, gin.H{"jwt": jwtValue})
}
