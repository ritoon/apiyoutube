package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	kafka "github.com/segmentio/kafka-go"

	"apiyoutube/queue"
)

type ServiceStats struct {
	kafkaReader *kafka.Reader
	cache       *redis.Client
}

func NewStats(cache *redis.Client, reader *kafka.Reader) *ServiceStats {

	ss := &ServiceStats{
		kafkaReader: reader,
		cache:       cache,
	}
	go ss.sendMessageToRedis()
	return ss
}

func (ss *ServiceStats) StatLogin(ctx *gin.Context) {
	now := fmt.Sprintf("%v-%v-%v", time.Now().Month, time.Now().Day, time.Now().Hour)
	res := ss.cache.PFCount(now)
	fmt.Printf("key %v res %v", now, res)
	ctx.JSON(http.StatusOK, gin.H{"stat": res.Val()})
}

func (ss *ServiceStats) sendMessageToRedis() {
	ctx := context.Background()
	for {
		m, err := ss.kafkaReader.FetchMessage(ctx)
		if err != nil {
			break
		}
		// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), ))
		var msg queue.MessageLogin
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Println(err)
		}

		now := fmt.Sprintf("%v-%v-%v", msg.EventTime.Month, msg.EventTime.Day, msg.EventTime.Hour)
		log.Println("get from kafka", m.Key, msg)
		pf := ss.cache.PFAdd(now, string(m.Key))
		if pf.Err != nil {
			fmt.Println("service/stats: err cache", pf.Err())
		}

		ss.kafkaReader.CommitMessages(ctx, m)
	}
}
