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
		Balancer: &kafka.RoundRobin{},
	})
	return &KafkaSender{Writer: writer}
}

func (kw *KafkaSender) SendMessage(ctx context.Context, key, value []byte) error {
	log.Printf("STRING: key: %s, value: %s", string(key), string(value))
	err := kw.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to write in topic: %v", err)
	}

	return nil
}
