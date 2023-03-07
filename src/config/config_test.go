package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigCreate(t *testing.T) {
	os.Setenv("DISCORD_TOKEN", "something")
	os.Setenv("OPEN_AI_API_KEY", "last")
	defer os.Unsetenv("DISCORD_TOKEN")
	defer os.Unsetenv("OPEN_AI_API_KEY")

	config := CreateConfig()

	assert.Equal(t, config.DiscordToken, "something")
	assert.Equal(t, config.OpenAIApiKey, "last")
}
