package componenets

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot-go-template/sos"
)

func FailedRescue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sosID := strings.Split(i.MessageComponentData().CustomID, ":")[1]
	call := sos.GetSos(sosID)

	if call.Status == sos.Open || call.Status == sos.Responded {
		call.Failed()

		// update the original sos message
		edit := discordgo.NewMessageEdit(call.ChannelID, call.MessageID)
		message, err := s.ChannelMessage(call.ChannelID, call.MessageID)
		if err != nil {
			panic(err)
		}

		b := message.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.Button)
		b.Disabled = true
		edit.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					b,
				},
			},
		}
		edit.Embeds = []*discordgo.MessageEmbed{
			{
				Title: "Rescue Information",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Who",
						Value: message.Embeds[0].Fields[0].Value,
					},
					{
						Name:  "Where",
						Value: message.Embeds[0].Fields[1].Value,
					},
					{
						Name:  "Responder",
						Value: message.Embeds[0].Fields[2].Value,
					},
					{
						Name:  "Status",
						Value: "Failed",
					},
				},
			},
		}

		_, err = s.ChannelMessageEditComplex(edit)
		if err != nil {
			panic(err)
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		}); err != nil {
			panic(err)
		}

		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	}); err != nil {
		panic(err)
	}
}
