package commands

import (
	"fmt"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot/bank"
)

var subCommands = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"balance":     balance,
	"transaction": transaction,
}

// Bank command handler
func Bank(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	if sc, ok := subCommands[options[0].Name]; ok {
		sc(s, i)
	}
}

// Balance sub command handler
func balance(s *discordgo.Session, i *discordgo.InteractionCreate) {
	b, err := bank.GetBank(i.GuildID)
	if err != nil {
		log.WithError(err).Error("balance")
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Check the org balance",
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("%daUEC", b.Balance),
		},
	}); err != nil {
		panic(err)
	}
}

// Transaction sub command handler
func transaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "bank-transaction",
			Title:    "Make a new transaction",
			Flags:    discordgo.MessageFlagsEphemeral,
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
