package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})
	return &Consumer{reader: r}
}

func (c *Consumer) Listen(ctx context.Context, handler func(key, value []byte)) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("error reading message: %v", err)
			continue
		}
		handler(msg.Key, msg.Value)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
