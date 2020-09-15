package main

import (
	"github.com/gin-gonic/gin"

	"apiyoutube/db/mock"
	"apiyoutube/middleware"
	"apiyoutube/service"
)

func main() {
	db := mock.New()
	su := service.NewUser(db)

	r := gin.Default()
	r.Use(middleware.VerifyJWT())
	r.GET("/user/:uuid", su.GetUser)
	r.GET("/user", su.GetListUser)
	r.POST("/user", su.CreateUser)
	r.Run()
}
