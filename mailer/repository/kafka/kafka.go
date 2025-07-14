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
		Brokers: brokers,
		Topic:   topicName,
		GroupID: groupID,
		// StartOffset:           kafka.FirstOffset,
		// MinBytes:              1,    // 1B
		// MaxBytes:              10e6, // 10MB
		CommitInterval:        time.Second,
		WatchPartitionChanges: true,
	})
	return &KafkaRecipient{Reader: reader}
}

func (kr *KafkaRecipient) ReceiveMessage(ctx context.Context) (key, value string, err error) {
	log.Println("reading message...")
	message, err := kr.Reader.ReadMessage(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to read message: %v", err)
	}

	key = string(message.Key)
	value = string(message.Value)
	log.Printf("message has been read: key: %s, value: %s\n", key, value)
	return key, value, nil
}
