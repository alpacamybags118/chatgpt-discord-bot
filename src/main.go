package main

import (
	chathandler "chatgpt-discord-bot/src/chat"
	"chatgpt-discord-bot/src/commands"
	"chatgpt-discord-bot/src/config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session
var configSettings *config.Config

func init() {
	var err error
	configSettings = config.CreateConfig()
	session, err = discordgo.New("Bot " + configSettings.DiscordToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func init() {
	commandHandlers := commands.GetCommandHandlers()
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		input := chathandler.CommandHandlerInput{
			Interaction: i,
			Session:     s,
			Config:      configSettings,
		}

		if h, ok := commandHandlers[input.Interaction.ApplicationCommandData().Name]; ok {
			h(input)
		}
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
		if i.Author.ID == s.State.User.ID {
			return
		}

		ch, err := s.State.Channel(i.ChannelID)

		if err != nil {
			log.Printf("Error getting channel: %s", err)
			return
		}

		if ch.IsThread() && ch.OwnerID == session.State.User.ID {
			input := chathandler.ReplyChatInput{
				Session: s,
				Config:  configSettings,
				Thread:  ch,
			}

			chathandler.ReplyInChat(input)
		}
	})
}

func main() {
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	session.Identify.Intents |= discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	err = commands.PushCommands(session, configSettings)
	if err != nil {
		log.Panicf("Cannot create commands: %v", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
