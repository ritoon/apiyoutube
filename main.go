package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"apiyoutube/db/mock"
	"apiyoutube/middleware"
	"apiyoutube/service"
)

func main() {

	conf := getconfig()

	db := mock.New()
	su := service.NewUser(db)

	// init router
	r := gin.Default()
	r.Use(middleware.VerifyJWT(conf.JWTSecret))
	r.GET("/user/:uuid", su.GetUser)
	r.GET("/user", su.GetListUser)
	r.POST("/user", su.CreateUser)
	r.Run()
}

type Config struct {
	JWTSecret string
}

func getconfig() *Config {
	// get config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	var conf Config
	conf.JWTSecret = viper.GetString("SECRET_KEY_JWT")
	return &conf
}
