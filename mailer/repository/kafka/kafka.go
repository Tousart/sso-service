package kafka

import (
	"context"
	"fmt"
	"log"

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
	})
	return &KafkaRecipient{Reader: reader}
}

func (kr *KafkaRecipient) ReceiveMessage(ctx context.Context) (key, value string, err error) {
	message, err := kr.Reader.ReadMessage(ctx)
	log.Printf("BYTES: key: %v, value: %v\n", key, value)
	if err != nil {
		return "", "", fmt.Errorf("failed to read message: %v", err)
	}

	key = string(message.Key)
	value = string(message.Value)
	log.Printf("STRING: key: %s, value: %s\n", key, value)
	return key, value, nil
}
