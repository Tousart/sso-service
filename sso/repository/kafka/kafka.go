package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaSender struct {
	Writer *kafka.Writer
}

func NewKafkaSender(brokers []string, topicName string) *KafkaSender {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topicName,
		Balancer: &kafka.Hash{},
	})
	return &KafkaSender{Writer: writer}
}

func (kw *KafkaSender) SendMessage(ctx context.Context, key, value []byte) error {
	log.Printf("sending message: key: %s, value: %s", string(key), string(value))
	err := kw.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)

	log.Println("message has been sent to broker")

	if err != nil {
		return fmt.Errorf("failed to write in topic: %v", err)
	}

	return nil
}
