package handler

import (
	"fmt"

	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/kafka"
)

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{Kafka: kafka}
}

func (h *TransactionCreatedKafkaHandler) Handle(message events.EventInterface) {
	if err := h.Kafka.Publish(message.GetPayload(), nil, "transactions"); err != nil {
		fmt.Println("TransactionCreatedKafkaHandler error:", err)
		return
	}
	fmt.Println("TransactionCreatedKafkaHandler published:", message.GetPayload())
}
