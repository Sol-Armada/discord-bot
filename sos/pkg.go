package sos

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/sol-armada/discord-bot-go-template/store"
)

type status string

const (
	Open      status = "open"
	Responded status = "responded"
	Rescued   status = "rescued"
	Failed    status = "failed"
	Canceled  status = "canceled"
)

type Sos struct {
	ID          string    `json:"id"`
	ChannelID   string    `json:"channel_id"`
	MessageID   string    `json:"message_id"`
	PlayerID    string    `json:"player_id"`
	MedicID     string    `json:"medic_id"`
	Where       string    `json:"where"`
	CallTime    time.Time `json:"call_time"`
	RespondTime time.Time `json:"respond_time"`
	RescuedTime time.Time `json:"rescued_time"`
	Status      status    `json:"status"`

	ctx context.Context
}

var calls map[string]*Sos = map[string]*Sos{}

func init() {
	ctx := context.Background()
	keys, err := store.Client.Keys(ctx, "sos:*").Result()
	if err != nil {
		log.WithError(err).Error("could not get list of stored sos calls")
	}

	for _, v := range keys {
		r, err := store.Client.Get(ctx, v).Result()
		if err != nil {
			log.WithError(err).Error("could not get stored sos call")
		}

		call := &Sos{}
		err = json.Unmarshal([]byte(r), call)
		if err != nil {
			log.WithError(err).Error("could not unmarshal stored sos call")
		}

		calls[call.ID] = call
	}
}

func New(i *discordgo.Interaction, where string) string {
	call := &Sos{
		ID:        xid.New().String(),
		ChannelID: i.ChannelID,
		PlayerID:  i.Member.User.ID,
		Where:     where,
		CallTime:  time.Now().UTC(),

		Status: Open,

		ctx: context.Background(),
	}
	call.StartCountDown()

	calls[call.ID] = call

	call.Store()

	return call.ID
}

func GetCalls() map[string]*Sos {
	return calls
}

func GetSos(id string) *Sos {
	return calls[id]
}

func (s *Sos) SetMessageID(id string) {
	s.MessageID = id

	s.Store()
}

func (s *Sos) StartCountDown() {
	go func() {
		// dead at 1.5 hours
		diedAt := s.CallTime.Add(90 * time.Minute)

		for {
			time.Sleep(10 * time.Second)
			if s.CallTime.After(diedAt) {
				s.Failed()
			}

			if s.Status == Canceled || s.Status == Failed || s.Status == Rescued {
				break
			}
		}
	}()
}

func (s *Sos) OnTheWay(m *discordgo.Member) {
	s.MedicID = m.User.ID
	s.RespondTime = time.Now().UTC()
	s.Status = Responded

	s.Store()
}

func (s *Sos) Canceled() {
	s.Status = Canceled

	s.Store()
}

func (s *Sos) Rescued() {
	s.Status = Rescued
	s.RescuedTime = time.Now().UTC()

	s.Store()
}

func (s *Sos) Failed() {
	s.Status = Failed

	s.Store()
}

func (s *Sos) Store() {
	key := fmt.Sprintf("sos:%s:%s", s.ID, s.PlayerID)

	m, err := json.Marshal(s)
	if err != nil {
		log.WithError(err).Error("could not store sos")
	}

	if err := store.Client.Set(s.ctx, key, string(m), 0).Err(); err != nil {
		log.WithError(err).Error("could not store sos")
	}
}
