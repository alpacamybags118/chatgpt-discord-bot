package commands

import (
	chathandler "chatgpt-discord-bot/src/chat"
	"chatgpt-discord-bot/src/config"

	"github.com/bwmarrin/discordgo"
)

func PushCommands(session *discordgo.Session, config *config.Config) error {
	commandsToPush := getCommands()

	for _, v := range commandsToPush {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, config.GuildID, v)
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
	}
}

func GetCommandHandlers() map[string]func(input chathandler.StartChatInput) {
	return map[string]func(input chathandler.StartChatInput){
		"start-chat": func(input chathandler.StartChatInput) {
			chathandler.StartChat(input)
		},
	}
}
