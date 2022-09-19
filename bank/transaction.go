package bank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sol-armada/discord-bot/store"
)

type Transaction struct {
	Bank           *Bank  `json:"bank"`
	ID             string `json:"id"`
	To             string `json:"to"`
	From           string `json:"from"`
	Amount         int64  `json:"ammount"`
	Status         status `json:"status"`
	RejectedReason string `json:"deny_reason"`

	ctx context.Context
}

func NewTransaction(bank *Bank, to string, from string, amount int64) (*Transaction, error) {
	t := &Transaction{
		Bank:   bank,
		ID:     xid.New().String(),
		To:     to,
		From:   from,
		Amount: amount,
		Status: Pending,
		ctx:    context.Background(),
	}
	if err := t.save(); err != nil {
		return nil, err
	}

	return t, nil
}

func GetTransactionById(guildId string, id string) (*Transaction, error) {
	key := fmt.Sprintf("%s:bank:transactions:%s", guildId, id)

	cmd := store.Client.Get(context.Background(), key)
	transaction, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	t := &Transaction{}
	if err := json.Unmarshal([]byte(transaction), t); err != nil {
		return nil, err
	}

	return t, nil
}

func GetTransactions() {}

func (t *Transaction) save() error {
	key := fmt.Sprintf("%s:bank:transactions:%s", t.Bank.Guild, t.ID)

	m, err := json.Marshal(t)
	if err != nil {
		return errors.Wrap(err, "unmarshal bank transaction")
	}

	if err := store.Client.Set(t.ctx, key, string(m), 0).Err(); err != nil {
		return errors.Wrap(err, "save bank transaction")
	}

	return nil
}

func (t *Transaction) Canceled() {
	t.Status = Canceled
	if err := t.save(); err != nil {
		log.WithError(err).Error("transaction rejected")
	}
}

func (t *Transaction) Rejected(reason string) {
	t.Status = Rejected
	t.RejectedReason = reason
	if err := t.save(); err != nil {
		log.WithError(err).Error("transaction rejected")
	}
}

func (t *Transaction) Processed() {
	if t.To == "bank" {
		t.Bank.AddBalance(t.Amount)
	} else {
		t.Bank.RemoveBalance(t.Amount)
	}

	t.Status = Processed
	if err := t.save(); err != nil {
		log.WithError(err).Error("proccessed")
	}
}
