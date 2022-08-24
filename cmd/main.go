package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot/bot"
	"github.com/sol-armada/discord-bot/settings"
	"github.com/sol-armada/discord-bot/starmap"
)

func main() {
	// get the settings
	settings.SetConfigName("settings")
	settings.AddConfigPath(".")

	err := settings.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := starmap.Load(); err != nil {
		panic(err)
	}

	appId := settings.GetString("APP_ID")

	//TODO: Get version and build
	log.WithFields(log.Fields{
		"app_id": appId,
	})

	// create a discrod session
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", settings.GetString("TOKEN")))
	if err != nil {
		panic(err)
	}

	b := bot.New(discord)

	// start the bot
	if err := b.Start(&bot.Options{
		AppID: appId,
	}); err != nil {
		log.WithError(err).Error("Failed to start the api")
		return
	}
}
