package sos

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
)

type Sos struct {
	ID          string
	Message     *discordgo.Message
	Who         *discordgo.Member
	Medic       *discordgo.Member
	Where       string
	CallTime    time.Time
	RespondTime time.Time
	RescuedTime time.Time
	GotRescued  bool
	Died        bool
}

var calls []*Sos = []*Sos{}

func New(i *discordgo.Interaction, where string) string {
	s := &Sos{
		ID:       xid.New().String(),
		Message:  i.Message,
		Who:      i.Member,
		Where:    where,
		CallTime: time.Now().UTC(),

		GotRescued: false,
		Died:       false,
	}
	s.StartCountDown()

	calls = append(calls, s)

	return s.ID
}

func GetCalls() []*Sos {
	return calls
}

func (s *Sos) StartCountDown() {
	go func() {
		deadAt := s.CallTime.Add(90 * time.Minute)

		for {
			if s.GotRescued {
				break
			}
			if s.CallTime.After(deadAt) {
				s.Died = true
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func (s *Sos) OnTheWay(m *discordgo.Member) {
	s.Medic = m
	s.RespondTime = time.Now().UTC()
}

func (s *Sos) Rescued() {
	s.GotRescued = true
	s.RescuedTime = time.Now().UTC()
}
