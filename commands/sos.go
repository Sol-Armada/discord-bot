package commands

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot/starmap"
)

// Sos command handler
func Sos(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "sos",
			Title:    "Call for help",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "location",
							Label:       "Specify a location",
							Style:       discordgo.TextInputShort,
							Placeholder: fmt.Sprintf("Example: %s", starmap.Keys[rand.Intn(len(starmap.Keys))]),
							Required:    true,
							MinLength:   3,
						},
					},
				},
			},
		},
	}); err != nil {
		panic(err)
	}
}
