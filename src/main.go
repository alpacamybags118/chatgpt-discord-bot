package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config := Config()

	dg, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
}
