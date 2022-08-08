package modals

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sol-armada/discord-bot-go-template/sos"
	"github.com/sol-armada/discord-bot-go-template/starmap"
)

func Sos(s *discordgo.Session, i *discordgo.InteractionCreate) {

	where := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	whereMatch := fuzzy.Find(where, starmap.Keys)

	if len(whereMatch) == 0 {
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   uint64(discordgo.MessageFlagsEphemeral),
				Content: fmt.Sprintf("I could not find a location named '%s'. Please try again", where),
			},
		}); err != nil {
			panic(err)
		}
		return
	}

	// create the sos call
	sosID := sos.New(i.Interaction, whereMatch[0])

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "sos",
			Content:  "Rescue needed!",
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
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: fmt.Sprintf("on-my-way-%s", sosID),
							Label:    "On my way!",
						},
					},
				},
			},
		},
	}); err != nil {
		panic(err)
	}

	_, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: "Click the button below if you no longer need a rescue or you died",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: fmt.Sprintf("cancel-rescue-%s", sosID),
						Label:    "Cancel",
					},
					discordgo.Button{
						CustomID: fmt.Sprintf("failed-rescue-%s", sosID),
						Label:    "Died",
					},
				},
			},
		},
		Username: i.Member.User.Username,
		Flags:    uint64(discordgo.MessageFlagsEphemeral),
	})
	if err != nil {
		panic(err)
	}
}
