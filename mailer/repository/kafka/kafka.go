package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaRecipient struct {
	Reader *kafka.Reader
}

func NewKafkaRecipient(brokers []string, topicName, groupID string) *KafkaRecipient {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               brokers,
		Topic:                 topicName,
		GroupID:               groupID,
		StartOffset:           kafka.FirstOffset,
		MinBytes:              1,
		MaxBytes:              10e6,
		MaxWait:               10 * time.Second,
		CommitInterval:        time.Second,
		WatchPartitionChanges: true,
	})
	return &KafkaRecipient{Reader: reader}
}

func (kr *KafkaRecipient) ReceiveMessage(ctx context.Context) (key string, value []byte, err error) {
	log.Println("reading message...")
	message, err := kr.Reader.ReadMessage(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read message: %v", err)
	}

	key = string(message.Key)
	value = message.Value
	log.Printf("message has been read: key: %s, value: %s\n", key, string(value))
	return key, value, nil
}
