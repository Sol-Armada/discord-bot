package bank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/sol-armada/discord-bot/store"
)

type Bank struct {
	Balance int64  `json:"balance"`
	Guild   string `json:"guild"`

	ctx context.Context
}

func GetBank(guildId string) (*Bank, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:bank", guildId)

	b := &Bank{
		ctx: context.Background(),
	}

	// check if we already have this guild's bank
	if !store.Exists(ctx, key) {
		b.Balance = 0
		b.Guild = guildId
		b.save()
		return b, nil
	}

	// get the existing bank
	cmd := store.Client.Get(ctx, key)
	a, err := cmd.Result()
	if err != nil {
		log.WithError(err).Error("getting bank from store")
		return nil, err
	}

	err = json.Unmarshal([]byte(a), b)
	if err != nil {
		log.WithError(err).Error("unmarshal existing bank")
	}

	return b, nil
}

func (b *Bank) AddBalance(amount int64) {
	b.Balance += amount
	b.save()
}

func (b *Bank) RemoveBalance(amount int64) {
	b.Balance -= amount
	b.save()
}

func (b *Bank) save() {
	key := fmt.Sprintf("%s:bank", b.Guild)

	m, err := json.Marshal(b)
	if err != nil {
		log.WithError(err).Error("could not save bank request")
	}

	if err := store.Client.Set(b.ctx, key, string(m), 0).Err(); err != nil {
		log.WithError(err).Error("could not save bank request")
	}
}
