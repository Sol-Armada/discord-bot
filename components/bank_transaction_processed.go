package components

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot/bank"
	"github.com/sol-armada/discord-bot/settings"
)

func BankTransactionProcessed(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// make sure they are allowed to do this
	if settings.IsInStringList("BANK.WHITE_LIST", i.Member.User.ID) {
		transactionId := strings.Split(i.MessageComponentData().CustomID, ":")[1]
		transaction, err := bank.GetTransactionById(i.GuildID, transactionId)
		if err != nil {
			log.WithError(err).Error("bank transaction processing component")
		}

		if transaction.Status == bank.Pending {
			transaction.Processed()

			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Marked trasaction as processed",
					Flags:   discordgo.MessageFlagsEphemeral,
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Transaction",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "From",
									Value: transaction.From,
								},
								{
									Name:  "To",
									Value: transaction.To,
								},
								{
									Name:  "Amount",
									Value: fmt.Sprintf("%d", transaction.Amount),
								},
								{
									Name:  "New Org Balance",
									Value: fmt.Sprintf("%d", transaction.Bank.Balance),
								},
							},
						},
					},
				},
			}); err != nil {
				log.WithError(err).Error("bank trasaction processed response")
				return
			}

			// edit the original message to disable the buttons
			edit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
			b := i.Message.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.Button)
			b.Disabled = true
			edit.Components = []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: fmt.Sprintf("bank-transaction-processed:%s", transactionId),
					Label:    "Processed",
					Disabled: true,
				},
				discordgo.Button{
					CustomID: fmt.Sprintf("bank-transaction-canceled:%s", transactionId),
					Label:    "Canceled",
					Disabled: true,
				},
			}

			if _, err := s.ChannelMessageEditComplex(edit); err != nil {
				log.WithError(err).Error("bank transaction edit original interaction response")
				return
			}
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		}); err != nil {
			log.WithError(err).Error("bank transaction end response")
		}

		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "You don't have permission to mark this transaction as processed",
		},
	}); err != nil {
		log.WithError(err).Error("bank transaction not allowed response")
	}
}
