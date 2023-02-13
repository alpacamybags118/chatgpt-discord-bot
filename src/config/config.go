package config

import "os"

type Config struct {
	OpenAIApiKey string
	DiscordToken string
	ChatGPTUrl   string
}

func CreateConfig() *Config {
	var config Config

	config.DiscordToken = os.Getenv("DISCORD_TOKEN")
	config.ChatGPTUrl = os.Getenv("CHATGPT_URL")
	config.OpenAIApiKey = os.Getenv("OPEN_AI_API_KEY")

	return &config
}
