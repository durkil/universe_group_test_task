package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"universe_group_test_task/pkg/kafka"
)

func main() {
	kafkaBroker := getEnv("KAFKA_BROKER", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "product-events")
	kafkaGroup := getEnv("KAFKA_GROUP", "notifications-group")

	consumer := kafka.NewConsumer([]string{kafkaBroker}, kafkaTopic, kafkaGroup)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("shutting down notifications service...")
		cancel()
	}()

	log.Printf("notifications service listening on topic %q", kafkaTopic)

	if err := consumer.Listen(ctx, func(key, value []byte) {
		log.Printf("received event: key=%s value=%s", string(key), string(value))
	}); err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	log.Println("notifications service stopped")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
