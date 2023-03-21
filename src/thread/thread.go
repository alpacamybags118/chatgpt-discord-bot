package thread

import (
	"chatgpt-discord-bot/src/config"
	"log"

	"github.com/bwmarrin/discordgo"
)

// IsThread will return whether the provided channel is a thread
func IsThread(channelId string, session *discordgo.Session, config *config.Config) bool {
	threads, err := session.GuildThreadsActive(config.GuildID)

	if err != nil {
		log.Fatalf("Error fetching threads: %s", err.Error())
		return false
	}

	for _, thread := range threads.Threads {
		if thread.ID == channelId {
			return true
		}
	}

	return false
}
