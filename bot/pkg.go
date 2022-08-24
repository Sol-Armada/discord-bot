package bot

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/sol-armada/discord-bot/bank"
	"github.com/sol-armada/discord-bot/commands"
	"github.com/sol-armada/discord-bot/components"
	"github.com/sol-armada/discord-bot/modals"
	"github.com/sol-armada/discord-bot/settings"
	"github.com/spf13/viper"
)

type Bot struct {
	SOSChannel string
	SOSMessage string

	BankChannel string
	BankMessage string

	sess *discordgo.Session
}

type Options struct {
	AppID string
}

func New(session *discordgo.Session) *Bot {
	return &Bot{
		sess: session,
	}
}

func (b *Bot) Start(o *Options) error {
	// available command handlers
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"sos":  commands.Sos,
		"bank": commands.Bank,
	}

	// available component handlers
	componentHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"on-my-way":       components.OnMyWay,
		"cancel-rescue":   components.CancelRescue,
		"failed-rescue":   components.FailedRescue,
		"cancel-response": components.CancelResponse,
	}

	// available modal hanlers
	modalHandlers := map[string]func(session *discordgo.Session, i *discordgo.InteractionCreate){
		"sos": modals.Sos,
	}

	b.sess.AddHandler(func(session *discordgo.Session, r *discordgo.Ready) {
		log.Info("Bot Ready")
	})

	b.sess.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		switch interaction.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
				h(session, interaction)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[interaction.MessageComponentData().CustomID]; ok {
				h(session, interaction)
			}

			// very custom components
			cid := interaction.MessageComponentData().CustomID
			if strings.HasPrefix(cid, "on-my-way") {
				componentHandlers["on-my-way"](session, interaction)
			}
			if strings.HasPrefix(cid, "cancel-rescue") {
				componentHandlers["cancel-rescue"](session, interaction)
			}
			if strings.HasPrefix(cid, "failed-rescue") {
				componentHandlers["failed-rescue"](session, interaction)
			}
			if strings.HasPrefix(cid, "cancel-response") {
				componentHandlers["cancel-response"](session, interaction)
			}
		case discordgo.InteractionModalSubmit:
			if h, ok := modalHandlers[interaction.ModalSubmitData().CustomID]; ok {
				h(session, interaction)
			}
		}
	})

	// sos
	if settings.GetBoolWithDefault("SOS.ENABLED", false) {
		b.sess.AddHandler(func(session *discordgo.Session, guild *discordgo.GuildCreate) {
			defer func() {
				log.Info("SOS enabled")
			}()

			log.Info("enabling SOS")
			done := make(chan error)

			go func() {
				// see if we have a channel ID to use
				channelID := viper.GetString("SOS.CHANNEL_ID")
				if channelID != "" {
					b.SOSChannel = channelID
					_, err := b.sess.Channel(channelID)
					if err != nil {
						done <- errors.Wrap(err, "configured channel not found")
						return
					}
					done <- nil
					return
				}

				// if we don't have a channel ID, check if there is
				// already a channel with the specified name
				channelName := settings.GetStringWithDefault("SOS.CHANNEL_NAME", "✚-emergency-sos")
				channels, err := b.sess.GuildChannels(guild.ID)
				if err != nil {
					done <- errors.Wrap(err, "could not get guild channels")
					return
				}
				for _, channel := range channels {
					if channel.Name == channelName {
						b.SOSChannel = channel.ID
						done <- nil
						return
					}
				}
				log.Debugf("could not find an existing channel named %s", channelName)

				// if we don't have anything, creat the channel
				c, err := session.GuildChannelCreate(guild.ID, channelName, discordgo.ChannelTypeGuildText)
				if err != nil {
					done <- errors.Wrap(err, "could not create SOS channel")
					return
				}
				b.SOSChannel = c.ID
				done <- nil
			}()

			err := <-done
			if err != nil {
				log.WithError(err).Error("could not enable SOS")
				b.SOSChannel = ""
				return
			}

			// register the commands with discord
			for _, v := range commands.SosCommands {
				_, err = b.sess.ApplicationCommandCreate(o.AppID, guild.ID, v)
				if err != nil {
					log.WithError(err).Error("could not enable SOS commands")
					b.SOSChannel = ""
				}
			}
		})
	}

	// bank
	if settings.GetBoolWithDefault("BANK.ENABLED", false) {
		b.sess.AddHandler(func(session *discordgo.Session, guild *discordgo.GuildCreate) {
			defer func() {
				log.Info("Bank enabled")
			}()

			log.Info("enabling Bank")
			done := make(chan error)
			defer close(done)

			go func() {
				// see if we have a channel ID to use
				channelID := viper.GetString("BANK.CHANNEL_ID")
				if channelID != "" {
					b.BankChannel = channelID
					_, err := b.sess.Channel(channelID)
					if err != nil {
						done <- errors.Wrap(err, "configured channel not found")
						return
					}
					done <- nil
					return
				}

				// if we don't have a channel ID, check if there is
				// already a channel with the specified name
				channelName := settings.GetStringWithDefault("BANK.CHANNEL_NAME", "💳-org-bank")
				channels, err := b.sess.GuildChannels(guild.ID)
				if err != nil {
					done <- errors.Wrap(err, "could not get guild channels")
					return
				}
				for _, channel := range channels {
					if channel.Name == channelName {
						b.SOSChannel = channel.ID
						done <- nil
						return
					}
				}
				log.Debugf("could not find an existing channel named %s", channelName)

				// if we don't have anything, creat the channel
				c, err := session.GuildChannelCreate(guild.ID, channelName, discordgo.ChannelTypeGuildText)
				if err != nil {
					done <- errors.Wrap(err, "could not create Bank channel")
					return
				}
				b.SOSChannel = c.ID
				done <- nil
			}()

			err := <-done
			if err != nil {
				log.WithError(err).Error("could not enable the Bank")
				b.SOSChannel = ""
				return
			}

			// register the commands with discord
			for _, v := range commands.Bankcommands {
				_, err = b.sess.ApplicationCommandCreate(o.AppID, guild.ID, v)
				if err != nil {
					log.WithError(err).Error("could not enable Bank commands")
					b.SOSChannel = ""
				}
			}

			// create the guild's bank
			bank.GetBank(guild.ID)
		})
	}

	// open the connection to discord
	if err := b.sess.Open(); err != nil {
		return err
	}
	defer b.sess.Close()

	// catch the kill signal for a clean shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill) // nolint:staticcheck
	<-stop
	log.Info("Bot is shutting down")

	// remove all available commands
	c, err := b.sess.ApplicationCommands(o.AppID, "")
	if err != nil {
		panic(err)
	}
	for _, v := range c {
		if err := b.sess.ApplicationCommandDelete(b.sess.State.User.ID, "", v.ID); err != nil {
			panic(err)
		}
	}

	return nil
}
