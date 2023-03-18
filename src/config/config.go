package config

import "os"

type Config struct {
	OpenAIApiKey     string
	DiscordToken     string
	GuildID          string
	DiscordPublicKey string
}

func CreateConfig() *Config {
	var config Config

	config.DiscordToken = os.Getenv("DISCORD_TOKEN")
	config.OpenAIApiKey = os.Getenv("OPEN_AI_API_KEY")
	config.GuildID = os.Getenv("GUILD_ID")
	config.DiscordPublicKey = os.Getenv("DISCORD_PUBLIC_KEY")

	return &config
}
