package handler

import (
	"fmt"

	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/events"
	"github.com.br/viniciusidacruz/microservice-wallet-core/pkg/kafka"
)

type BalanceUpdatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewBalanceUpdatedKafkaHandler(kafka *kafka.Producer) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{Kafka: kafka}
}

func (h *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface) {
	eventMessage := map[string]interface{}{
		"name":    message.GetName(),
		"payload": message.GetPayload(),
	}

	if err := h.Kafka.Publish(eventMessage, nil, "balance_updated"); err != nil {
		fmt.Println("BalanceUpdatedKafkaHandler error:", err)
		return
	}

	fmt.Println("BalanceUpdatedKafkaHandler published:", eventMessage)
}
