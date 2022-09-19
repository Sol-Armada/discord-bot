package modals

import (
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot/bank"
)

func BankTo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// get the amount
	amount, err := strconv.ParseInt(i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value, 10, 64)
	if err != nil || amount <= 0 {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "That is not a correct amount. Must be a valid number greater than 0",
			},
		}); err != nil {
			log.WithError(err).Error("bank transction to modal interaction resonse")
		}
		return
	}

	b, err := bank.GetBank(i.GuildID)
	if err != nil {
		log.WithError(err).Error("get bank")
	}

	t, err := bank.NewTransaction(b, "bank", i.Member.User.ID, amount)
	if err != nil {
		log.WithError(err).Error("new bank transaction")
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Your transction has been submitted for processing",
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Amount",
							Value: fmt.Sprintf("%d", amount),
						},
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: fmt.Sprintf("bank-transaction-canceled:%s", t.ID),
							Label:    "Canceled",
						},
					},
				},
			},
		},
	}); err != nil {
		log.WithError(err).Error("bank transaction to modal interaction resonse")
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s is transfering %daUEC to the org bank", i.Member.User.Mention(), amount),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: fmt.Sprintf("bank-transaction-processed:%s", t.ID),
						Label:    "Processed",
					},
					discordgo.Button{
						CustomID: fmt.Sprintf("bank-transaction-canceled:%s", t.ID),
						Label:    "Canceled",
					},
				},
			},
		},
	}); err != nil {
		log.WithError(err).Error("bank transcation to follow up message")
	}
}
