package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Sos command handler
func Sos(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "test",
			Content: "TEST",
		},
	}); err != nil {
		panic(err)
	}
}
