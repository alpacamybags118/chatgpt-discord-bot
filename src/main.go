package main

import (
	"chatgpt-discord-bot/src/config"
	opengptclient "chatgpt-discord-bot/src/opengpt"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config := config.CreateConfig()

	dg, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	defer dg.Close()

	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	config := config.CreateConfig()
	client := opengptclient.CreateNew(&config)
	req := opengptclient.OpenGptCompletionRequest{
		Prompt:      m.Content,
		Temperature: 1.0,
		Model:       "text-davinci-003",
		Max_tokens:  200,
	}

	resp, err := client.SendCompletionRequest(req)

	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "error occured")
	}

	s.ChannelMessageSend(m.ChannelID, resp.Choices[0].Text)
}
