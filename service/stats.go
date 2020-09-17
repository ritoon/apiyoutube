package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"

	"apiyoutube/queue"
)

type ServiceStats struct {
	kafkaReader func(string, string) *kafka.Reader
	cache       *redis.Client
	elastic     *elasticsearch.Client
}

func NewStats(cache *redis.Client, elastic *elasticsearch.Client, reader func(string, string) *kafka.Reader) *ServiceStats {

	ss := &ServiceStats{
		kafkaReader: reader,
		cache:       cache,
		elastic:     elastic,
	}
	go ss.sendMessageToRedis()
	go ss.sendMessageToElastic()
	return ss
}

func (ss *ServiceStats) StatLogin(ctx *gin.Context) {
	now := fmt.Sprintf("%v-%v-%v", time.Now().Month, time.Now().Day, time.Now().Hour)
	res := ss.cache.PFCount(now)
	fmt.Printf("key %v res %v", now, res)
	ctx.JSON(http.StatusOK, gin.H{"stat": res.Val()})
}

func (ss *ServiceStats) sendMessageToElastic() {
	reader := ss.kafkaReader("group-elastic", queue.TopicUserLogin)
	ctx := context.Background()
	for {
		msg, origin, err := ss.parseMessage(ctx, reader)
		if err != nil {
			log.Println(err)
		}
		log.Println(msg, origin)

		req := esapi.IndexRequest{
			Index:      "userlogin",
			DocumentID: uuid.New().String(),
			Body:       bytes.NewBuffer(origin.Value),
			Refresh:    "true",
		}

		// Perform the request with the client.
		res, err := req.Do(context.Background(), ss.elastic)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		log.Println(res)
		reader.CommitMessages(ctx, *origin)
	}
}

func (ss *ServiceStats) sendMessageToRedis() {
	reader := ss.kafkaReader("group-redis", queue.TopicUserLogin)
	ctx := context.Background()
	for {

		msg, origin, err := ss.parseMessage(ctx, reader)
		if err != nil {
			log.Println(err)
		}
		now := fmt.Sprintf("%v-%v-%v", msg.EventTime.Month, msg.EventTime.Day, msg.EventTime.Hour)
		pf := ss.cache.PFAdd(now, string(origin.Key))
		if pf.Err != nil {
			fmt.Println("service/stats: err cache", pf.Err())
		}
		reader.CommitMessages(ctx, *origin)
	}
}

func (ss *ServiceStats) parseMessage(ctx context.Context, reader *kafka.Reader) (*queue.MessageLogin, *kafka.Message, error) {
	m, err := reader.FetchMessage(ctx)
	if err != nil {
		return nil, nil, err
	}
	var msg queue.MessageLogin
	return &msg, &m, json.Unmarshal(m.Value, &msg)
}
