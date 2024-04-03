package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type NewProducerKafkaInit interface {
	ProduceMessages(topic string, message []byte) error
	Close() error
}

type Producer struct {
	writer *kafka.Writer
}

func ProducerCafkaInit(brokers []string) (NewProducerKafkaInit, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: true,
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) ProduceMessages(topic string, message []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic:     topic,
		Value:     message,
		Partition: 0,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
