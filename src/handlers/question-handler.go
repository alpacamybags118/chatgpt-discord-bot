package handlers

import (
	"chatgpt-discord-bot/src/config"
	opengptclient "chatgpt-discord-bot/src/opengpt"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func HandleQuestion(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.ChannelMessageSend(i.ChannelID, "Asking your question...")

	options := i.ApplicationCommandData().Options
	fmt.Println("options")
	config := config.CreateConfig()
	client := opengptclient.CreateNew(config)
	req := opengptclient.OpenGptCompletionRequest{
		Prompt:      options[0].StringValue(),
		Temperature: 1.0,
		Model:       "text-davinci-003",
		Max_tokens:  200,
	}

	resp, err := client.SendCompletionRequest(req)

	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(i.ChannelID, "error occured")
	}

	s.ChannelMessageSend(i.ChannelID, resp.Choices[0].Text)
}
