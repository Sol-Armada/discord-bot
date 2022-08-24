package commands

import "github.com/bwmarrin/discordgo"

var SosCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "sos",
		Description: "Help!",
	},
}

var Bankcommands = []*discordgo.ApplicationCommand{
	{
		Name:        "bank",
		Description: "Money in. Money out.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "balance",
				Description: "Get the org balance",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
}
