package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ConsumerKafka interface {
	ConsumerMessages(handler func(message []byte)) error
	Close() error
}

type Consumer struct {
	reader *kafka.Reader
}

func ConsumerCafkaInit(brokers []string, topic , groupID string) (ConsumerKafka, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Consumer{reader: reader}, nil
}

func (c *Consumer) ConsumerMessages(handler func(message []byte)) error {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		handler(msg.Value)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}