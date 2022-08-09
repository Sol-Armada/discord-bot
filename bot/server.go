package bot

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/sol-armada/discord-bot-go-template/commands"
	componenets "github.com/sol-armada/discord-bot-go-template/components"
	"github.com/sol-armada/discord-bot-go-template/modals"
	"github.com/sol-armada/discord-bot-go-template/settings"
	"github.com/spf13/viper"
)

type Server struct {
	SOSChannel string
	SOSMessage string

	Sess   *discordgo.Session
	Logger *log.Entry
}

type Options struct {
	AppID string
}

func (s *Server) Start(o *Options) error {
	logger := s.Logger
	sess := s.Sess

	// register the available commands
	commands := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"sos": commands.Sos,
	}

	// register the available components
	components := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"on-my-way":       componenets.OnMyWay,
		"cancel-rescue":   componenets.CancelRescue,
		"failed-rescue":   componenets.FailedRescue,
		"cancel-response": componenets.CancelResponse,
	}

	// register the available modals
	modals := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"sos": modals.Sos,
	}

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info("Bot Ready")
	})

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commands[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := components[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}

			// very custom components
			cid := i.MessageComponentData().CustomID
			if strings.HasPrefix(cid, "on-my-way") {
				components["on-my-way"](s, i)
			}
			if strings.HasPrefix(cid, "cancel-rescue") {
				components["cancel-rescue"](s, i)
			}
			if strings.HasPrefix(cid, "failed-rescue") {
				components["failed-rescue"](s, i)
			}
			if strings.HasPrefix(cid, "cancel-response") {
				components["cancel-response"](s, i)
			}
		case discordgo.InteractionModalSubmit:
			if h, ok := modals[i.ModalSubmitData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	sess.AddHandler(func(session *discordgo.Session, guild *discordgo.GuildCreate) {
		defer func() {
			s.Logger.Info("SOS enabled")
		}()

		if settings.GetBoolWithDefault("SOS.ENABLED", false) {
			s.Logger.Info("enabling SOS")
			done := make(chan error)

			go func() {
				// see if we have a channel ID to use
				channelID := viper.GetString("SOS.CHANNEL_ID")
				if channelID != "" {
					s.SOSChannel = channelID
					_, err := sess.Channel(channelID)
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
				channels, err := sess.GuildChannels(guild.ID)
				if err != nil {
					done <- errors.Wrap(err, "could not get guild channels")
					return
				}
				for _, channel := range channels {
					if channel.Name == channelName {
						s.SOSChannel = channel.ID
						done <- nil
						return
					}
				}
				s.Logger.Debugf("could not find an existing channel named %s", channelName)

				// if we don't have anything, creat the channel
				c, err := session.GuildChannelCreate(guild.ID, channelName, discordgo.ChannelTypeGuildText)
				if err != nil {
					done <- errors.Wrap(err, "could not create SOS channel")
					return
				}
				s.SOSChannel = c.ID
				done <- nil
			}()

			err := <-done
			if err != nil {
				s.Logger.WithError(err).Error("could not enable SOS")
				s.SOSChannel = ""
				return
			}

			// register the command with discord
			_, err = sess.ApplicationCommandCreate(o.AppID, guild.ID, &discordgo.ApplicationCommand{
				Name:        "sos",
				Description: "HELP!",
			})
			if err != nil {
				s.Logger.WithError(err).Error("could not enable SOS command")
				s.SOSChannel = ""
			}
		}
	})

	// open the connection to discord
	if err := sess.Open(); err != nil {
		return err
	}
	defer sess.Close()

	// catch the kill signal for a clean shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill) // nolint:staticcheck
	<-stop
	logger.Info("Bot is shutting down")

	// remove all available commands
	c, err := sess.ApplicationCommands(o.AppID, "")
	if err != nil {
		panic(err)
	}
	for _, v := range c {
		if err := sess.ApplicationCommandDelete(sess.State.User.ID, "", v.ID); err != nil {
			panic(err)
		}
	}

	return nil
}
