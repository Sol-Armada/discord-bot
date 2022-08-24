package bank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/rs/xid"
	"github.com/sol-armada/discord-bot/store"
)

type status string

const (
	Transaction_PENDING   = "pending"
	Transaction_ACCEPTED  = "accepted"
	Transaction_PROCESSED = "processed"
	Transaction_REJECTED  = "rejected"
	Transaction_ERRORED   = "errored"
)

type Transaction struct {
	Bank       *Bank   `json:"bank"`
	ID         string  `json:"id"`
	To         string  `json:"to"`
	From       string  `json:"from"`
	Amount     float64 `json:"ammount"`
	Status     status  `json:"status"`
	DenyReason string  `json:"deny_reason"`

	ctx context.Context
}

func NewTransaction(bank *Bank, to string, from string, amount float64) *Transaction {
	t := &Transaction{
		ID:     xid.New().String(),
		To:     to,
		From:   from,
		Amount: amount,
		ctx:    context.Background(),
	}
	t.save()

	return t
}

func GetTransactionFrom() {}

func GetTransactionTo() {}

func GetTransactions() {}

func (t *Transaction) save() {
	key := fmt.Sprintf("%s:bank:transactions:%s", t.Bank.Guild, t.ID)

	m, err := json.Marshal(t)
	if err != nil {
		log.WithError(err).Error("unmarshal bank transaction")
	}

	if err := store.Client.Set(t.ctx, key, string(m), 0).Err(); err != nil {
		log.WithError(err).Error("save bank transaction")
	}
}
