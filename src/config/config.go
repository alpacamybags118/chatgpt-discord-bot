package config

type Config struct {
	OpenAIApiKey string `env:"OPEN_AI_API_KEY"`
	DiscordToken string `env:"DISCORD_TOKEN"`
}
