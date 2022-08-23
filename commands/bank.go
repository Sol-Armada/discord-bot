package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Balance command handler
func Balance(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "bank-balance",
			Title:    "Check the org balance",
			Flags:    uint64(discordgo.MessageFlagsEphemeral),
			Content:  "Future balance here",
		},
	}); err != nil {
		panic(err)
	}
}

// Transaction command handler
func Transaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "bank-transaction",
			Title:    "Make a new transaction",
			Flags:    uint64(discordgo.MessageFlagsEphemeral),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "To Org Bank",
						},
						discordgo.Button{
							Label: "From Org Bank",
						},
					},
				},
			},
		},
	}); err != nil {
		panic(err)
	}
}
