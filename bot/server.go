package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/sol-armada/discord-bot-go-template/commands"
)

type Server struct {
	Sess *discordgo.Session

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
	components := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

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
		}
	})

	// register the command with discord
	_, err := sess.ApplicationCommandCreate(o.AppID, "", &discordgo.ApplicationCommand{
		Name:        "sos",
		Description: "HELP!",
	})
	if err != nil {
		return errors.Wrap(err, "sos command")
	}

	// open the connection to discord
	if err := sess.Open(); err != nil {
		return err
	}
	defer sess.Close()

	// catch the kill signal for a clean shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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
