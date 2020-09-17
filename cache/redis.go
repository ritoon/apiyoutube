package cache

import (
	"fmt"

	redis "github.com/go-redis/redis"
)

func New() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong := rdb.Ping()
	fmt.Println(pong)
	return rdb
}
