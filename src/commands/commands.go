package main

import (
	"chatgpt-discord-bot/src/config"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const APP_ID string = "1068334708057981022"

func RemoveCommands(session *discordgo.Session, config *config.Config) error {
	commands, err := session.ApplicationCommands(APP_ID, config.GuildID)

	if err != nil {
		log.Fatalf("Error fetching commands: %s", err.Error())
	}

	for _, v := range commands {
		err := session.ApplicationCommandDelete(APP_ID, config.GuildID, v.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func PushCommands(session *discordgo.Session, config *config.Config) error {
	commandsToPush := getCommands()

	for _, v := range commandsToPush {
		_, err := session.ApplicationCommandCreate(APP_ID, config.GuildID, v)
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
	config := config.CreateConfig()

	session, err := discordgo.New(fmt.Sprintf("Bot %s", config.DiscordToken))

	if err != nil {
		log.Fatal(err.Error())
	}

	PushCommands(session, config)
}
