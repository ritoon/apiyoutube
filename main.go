package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"apiyoutube/cache"
	"apiyoutube/db/orm"
	"apiyoutube/middleware"
	"apiyoutube/queue"
	"apiyoutube/service"
)

func main() {
	// get the current config
	conf := getconfig()
	// init app
	initApp(conf)
}

func initApp(conf *Config) {
	// create cache
	cache := cache.New()

	// create the queue.
	queueWriter, queueReader := queue.New(conf.KafkaHost)

	db := orm.New(conf.DBHost, conf.DBUser, conf.DBPass, conf.DBName, conf.DBPort)
	su := service.NewUser(db, cache, queueWriter, conf.JWTSecret)
	ss := service.NewStats(cache, queueReader)
	// init router.
	r := gin.Default()
	r.GET("/stats-login", ss.StatLogin)
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
	KafkaHost string // KAFKA_HOST: 127.0.0.1
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
	conf.KafkaHost = viper.GetString("KAFKA_HOST")
	return &conf
}
