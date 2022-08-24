package bank

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	"github.com/rs/xid"
	"github.com/sol-armada/discord-bot/store"
)

type RequestStatus string

const (
	Request_PENDING   RequestStatus = "pending"
	Request_ACCEPTED  RequestStatus = "accepted"
	Request_PROCESSED RequestStatus = "processed"
	Request_DENIED    RequestStatus = "denied"
)

type Request struct {
	Bank        *Bank        `json:"bank"`
	ID          string       `json:"id"`
	Amount      float64      `json:"amount"`
	SubmittedBy string       `json:"submitted_by"`
	AcceptedBy  string       `json:"accepted_by"`
	Reason      string       `json:"reason"`
	DenyReason  string       `json:"deny_reason"`
	Transaction *Transaction `json:"transaction"`

	ctx context.Context
}

func NewRequest(amount float64, subbmittedBy string, acceptedBy string, reason string) *Request {
	r := &Request{
		ID:          xid.New().String(),
		Amount:      amount,
		SubmittedBy: subbmittedBy,
		AcceptedBy:  acceptedBy,
		Reason:      reason,
		ctx:         context.Background(),
	}
	r.save()

	return r
}

func GetRequestById(id string) {}

func (r *Request) save() {
	key := fmt.Sprintf("bank:request:%s", r.ID)

	m, err := json.Marshal(r)
	if err != nil {
		log.WithError(err).Error("could not save bank request")
	}

	if err := store.Client.Set(r.ctx, key, string(m), 0).Err(); err != nil {
		log.WithError(err).Error("could not save bank request")
	}
}

func (r *Request) Accept() {}

func (r *Request) Deny(reason string) {}
