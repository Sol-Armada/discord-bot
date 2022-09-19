package components

import (
	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
)

func BankTo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "bank-to-modal",
			Title:    "Transer to the Org Bank",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "ammount",
							Label:       "How much are you transfering to the Org bank?",
							Style:       discordgo.TextInputShort,
							Placeholder: "0",
							Required:    true,
							MinLength:   1,
						},
					},
				},
			},
		},
	}); err != nil {
		log.WithError(err).Error("bank transaction to button response")
	}
}
