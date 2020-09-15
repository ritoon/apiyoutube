package main

import (
	"apiyoutube/db/mock"
	"apiyoutube/middleware"
	"apiyoutube/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := mock.New()
	su := service.NewUser(db)

	r := gin.Default()
	r.Use(middleware.VerifyJWT())
	r.GET("/user/:uuid", su.GetUser)
	r.Run()
}
