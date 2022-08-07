package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/sol-armada/discord-bot-go-template/bot"
	"github.com/sol-armada/discord-bot-go-template/starmap"
	"github.com/spf13/viper"
)

func main() {
	// get the settings
	viper.SetConfigName("settings")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := starmap.Load(); err != nil {
		panic(err)
	}

	appId := viper.GetString("APP_ID")

	//TODO: Get version and build
	logger := log.WithFields(log.Fields{
		"app_id": appId,
	})

	// create a discrod session
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", viper.GetString("TOKEN")))
	if err != nil {
		panic(err)
	}

	// create the bot
	b := &bot.Server{
		Sess: discord,
	}
	b.Logger = logger

	// start the bot
	if err := b.Start(&bot.Options{
		AppID: appId,
	}); err != nil {
		logger.WithError(err).Error("Failed to start the api")
		return
	}
}
