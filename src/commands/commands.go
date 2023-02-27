package commands

import (
	"chatgpt-discord-bot/src/config"
	"chatgpt-discord-bot/src/handlers"

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
			Name:        "ask-a-question",
			Description: "Ask OpenAI a question!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "The question you want to ask",
					Required:    true,
				},
			},
		},
	}
}

func GetCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ask-a-question": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			handlers.HandleQuestion(s, i)
		},
	}
}
