package kafka

import (
	"testing"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestNewKafkaProducer(t *testing.T) {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}
	producer := NewKafkaProducer(configMap)
	assert.NotNil(t, producer)
}

func TestPublish(t *testing.T) {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}
	producer := NewKafkaProducer(configMap)
	assert.NotNil(t, producer)
}
