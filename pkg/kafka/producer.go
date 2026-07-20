package kafka

import (
	"encoding/json"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(config *ckafka.ConfigMap) *Producer {
	return &Producer{ConfigMap: config}
}

func (p *Producer) Publish(msg interface{}, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.ConfigMap)
	if err != nil {
		return err
	}
	defer producer.Close()

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Key:            key,
		Value:          payload,
	}

	err = producer.Produce(message, nil)
	if err != nil {
		return err
	}

	if remaining := producer.Flush(15 * 1000); remaining > 0 {
		return fmt.Errorf("failed to flush kafka messages: %d remaining", remaining)
	}

	return nil
}
