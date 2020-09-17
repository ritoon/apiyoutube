package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"apiyoutube/db/orm"
	"apiyoutube/middleware"
	"apiyoutube/service"
)

func main() {
	// get the current config
	conf := getconfig()
	// init app
	initApp(conf)
}

func initApp(conf *Config) {
	// create tools.
	// db := mock.New()
	db := orm.New(conf.DBHost, conf.DBUser, conf.DBPass, conf.DBName, conf.DBPort)
	su := service.NewUser(db, conf.JWTSecret)
	// init router.
	r := gin.Default()
	r.POST("/login", su.LoginUser)
	user := r.Group("/user")
	user.Use(middleware.VerifyJWT(conf.JWTSecret))
	user.GET("/:uuid", su.GetUser)
	user.GET("", su.GetListUser)
	user.POST("", su.CreateUser)
	user.DELETE("/:uuid", su.DeleteUser)
	user.PUT("/:uuid", su.UpdateUser)
	r.Run()
}

type Config struct {
	JWTSecret string // SECRET_KEY_JWT: my_secret_key
	DBName    string // POSTGRES_DB: apiyoutube
	DBUser    string // POSTGRES_USER: apiyoutube
	DBPass    string // POSTGRES_PASSWORD: password
	DBPort    string // POSTGRES_PORT: 5432
	DBHost    string // POSTGRES_HOST: 127.0.0.1
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
	conf.DBName = viper.GetString("DB.POSTGRES_DB")
	conf.DBUser = viper.GetString("DB.POSTGRES_USER")
	conf.DBPass = viper.GetString("DB.POSTGRES_PASSWORD")
	conf.DBPort = viper.GetString("DB.POSTGRES_PORT")
	conf.DBHost = viper.GetString("DB.POSTGRES_HOST")
	return &conf
}
