package components

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot/sos"
)

func SosCancelResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sosID := strings.Split(i.MessageComponentData().CustomID, ":")[1]
	call := sos.GetSos(sosID)

	if call.Status == sos.Open || call.Status == sos.Responded {
		call.ClearResponder()

		edit := discordgo.NewMessageEdit(call.ChannelID, call.MessageID)
		message, err := s.ChannelMessage(call.ChannelID, call.MessageID)
		if err != nil {
			panic(err)
		}

		b := message.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.Button)
		b.Disabled = false
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
						Value: "No reponder yet",
					},
					{
						Name:  "Status",
						Value: "Open",
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
	}
}
