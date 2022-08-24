package bank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/sol-armada/discord-bot/store"
)

type Bank struct {
	Balance int64
	Guild   string

	ctx context.Context
}

func GetBank(guild string) (*Bank, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:bank", guild)

	b := &Bank{}

	// check if we already have this guild's bank
	if !store.Exists(ctx, key) {
		b = &Bank{
			Balance: 0,
			Guild:   guild,
			ctx:     context.Background(),
		}
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

func (b *Bank) GetBalance() float64 {
	key := fmt.Sprintf("%s:bank", b.Guild)

	cmd := store.Client.Get(b.ctx, key)
	_ = cmd
	return 0
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
