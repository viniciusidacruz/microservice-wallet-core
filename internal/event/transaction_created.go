package event

import "time"

type TransactionCreated struct {
	Name    string
	Payload interface{}
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{
		Name: "TransactionCreated",
	}
}

func (t *TransactionCreated) GetName() string {
	return t.Name
}

func (t *TransactionCreated) GetPayload() interface{} {
	return t.Payload
}

func (t *TransactionCreated) SetPayload(payload interface{}) {
	t.Payload = payload
}

func (t *TransactionCreated) GetDataTime() time.Time {
	return time.Now()
}
