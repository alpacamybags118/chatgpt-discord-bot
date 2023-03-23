package main

import (
	"chatgpt-discord-bot/src/config"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func RemoveCommands(session *discordgo.Session, config *config.Config) error {
	commands, err := session.ApplicationCommands(config.BotUserID, config.GuildID)

	if err != nil {
		log.Fatalf("Error fetching commands: %s", err.Error())
	}

	for _, v := range commands {
		err := session.ApplicationCommandDelete(config.BotUserID, config.GuildID, v.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func PushCommands(session *discordgo.Session, config *config.Config) error {
	commandsToPush := getCommands()

	for _, v := range commandsToPush {
		_, err := session.ApplicationCommandCreate(config.BotUserID, config.GuildID, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "start-chat",
			Description: "Start a chat with ChatGPT",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "The initial prompt to begin the chat",
					Required:    true,
				},
			},
		},
		{
			Name:        "reply",
			Description: "Reply to your ChatGPT chat session",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reply",
					Description: "The reply to the last chat message from ChatGPT",
					Required:    true,
				},
			},
		},
	}
}

func main() {
	var action = flag.String("action", "create", "Create or delete commands")

	flag.Parse()

	config := config.CreateConfig()

	session, err := discordgo.New(fmt.Sprintf("Bot %s", config.DiscordToken))

	if err != nil {
		log.Fatal(err.Error())
	}

	if strings.ToLower(*action) == "create" {
		PushCommands(session, config)
	} else {
		RemoveCommands(session, config)
	}

}
