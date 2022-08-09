package componenets

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot-go-template/sos"
)

func OnMyWay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sosID := strings.Split(i.MessageComponentData().CustomID, ":")[1]
	sos := sos.GetSos(sosID)
	// can't response to you own call
	if sos.PlayerID == i.Member.User.ID {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You can't respond to your own call!",
				Flags:   uint64(discordgo.MessageFlagsEphemeral),
			},
		}); err != nil {
			panic(err)
		}

		return
	}

	sos.OnTheWay(i.Member)

	edit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
	b := i.Message.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.Button)
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
					Value: i.Message.Embeds[0].Fields[0].Value,
				},
				{
					Name:  "Where",
					Value: i.Message.Embeds[0].Fields[1].Value,
				},
				{
					Name:  "Responder",
					Value: i.Member.Mention(),
				},
				{
					Name:  "Status",
					Value: "Responded",
				},
			},
		},
	}

	_, err := s.ChannelMessageEditComplex(edit)
	if err != nil {
		panic(err)
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	}); err != nil {
		panic(err)
	}

	_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Flags: uint64(discordgo.MessageFlagsEphemeral),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: fmt.Sprintf("cancel-response:%s", sosID),
						Label:    "Cancel",
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}
