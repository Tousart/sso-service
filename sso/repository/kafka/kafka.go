package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaSender struct {
	writer *kafka.Writer
}

func NewKafkaSender(brokers []string, topicName string) *KafkaSender {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topicName,
		Balancer: &kafka.RoundRobin{},
	})
	return &KafkaSender{writer: writer}
}

func (kw *KafkaSender) SendMessage(ctx context.Context, key, value []byte) error {
	err := kw.writer.WriteMessages(ctx,
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
