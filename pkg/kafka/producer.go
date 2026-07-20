package kafka

import (
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *ckafka.Producer
}

func NewKafkaProducer(config *ckafka.ConfigMap) *Producer {
	producer, err := ckafka.NewProducer(config)
	if err != nil {
		panic(err)
	}

	return &Producer{producer: producer}
}

func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	deliveryChan := make(chan ckafka.Event, 1)
	err = p.producer.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Key:            key,
		Value:          payload,
	}, deliveryChan)
	if err != nil {
		return err
	}

	event := <-deliveryChan
	message := event.(*ckafka.Message)
	if message.TopicPartition.Error != nil {
		return fmt.Errorf("failed to deliver kafka message: %w", message.TopicPartition.Error)
	}

	return nil
}

func (p *Producer) Close() {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
}
