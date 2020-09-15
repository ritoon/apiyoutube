package main

import (
	"github.com/gin-gonic/gin"

	"github.com/ritoon/cours/discover/apiyoutube/db/mock"
	"github.com/ritoon/cours/discover/apiyoutube/service"
)

func main() {
	db := mock.New()
	su := service.NewUser(db)

	r := gin.Default()
	r.Use(middleware.VerifyJWT())
	r.GET("/user/:uuid", su.GetUser)
	r.Run()
}

