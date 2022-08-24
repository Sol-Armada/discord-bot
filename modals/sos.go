package modals

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sol-armada/discord-bot/sos"
	"github.com/sol-armada/discord-bot/starmap"
)

func Sos(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// check that the location is a known one
	where := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	whereMatch := fuzzy.FindFold(where, starmap.Keys)

	if len(whereMatch) == 0 {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("I could not find a location named '%s'. Please try again", where),
			},
		}); err != nil {
			panic(err)
		}
		return
	}

	// create the sos call
	call := sos.New(i.Interaction, whereMatch[0])

	m, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: "Rescue needed!",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Rescue Information",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Who",
						Value: i.Member.Mention(),
					},
					{
						Name:  "Where",
						Value: whereMatch[0],
					},
					{
						Name:  "Responder",
						Value: "No responder yet",
					},
					{
						Name:  "Status",
						Value: "Open",
					},
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: fmt.Sprintf("on-my-way:%s", call.ID),
						Label:    "✚ On my way! ✚",
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	call.SetMessageID(m.ID)

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: fmt.Sprintf("cancel-rescue:%s", call.ID),
							Label:    "Cancel",
						},
						discordgo.Button{
							CustomID: fmt.Sprintf("failed-rescue:%s", call.ID),
							Label:    "Died",
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		panic(err)
	}
}
