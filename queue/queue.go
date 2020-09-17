package queue

import (
	"time"

	kafka "github.com/segmentio/kafka-go"
)

const (
	topic = "user_login"
)

type MessageLogin struct {
	UUID        string    `json:"uuid"`
	UserAgent   string    `json:"user_agent"`
	CallContext string    `json:"call_context"`
	IP          string    `json:"ip"`
	EventTime   time.Time `json:"event_time"`
}

func New(host string) (*kafka.Writer, *kafka.Reader) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{host},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{host}, //"localhost:9092"}
		GroupID:        "consumer-group-id",
		Topic:          topic,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
		MaxWait:        time.Second / 10,
	})
	return w, r
}
