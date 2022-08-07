package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Sos command handler
func Sos(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal_test",
			Title:    "test",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "test",
							Label:       "test?",
							Style:       discordgo.TextInputShort,
							Placeholder: "test",
							Required:    true,
							MaxLength:   300,
							MinLength:   10,
						},
					},
				},
			},
		},
	}); err != nil {
		panic(err)
	}
}
