package chat

import (
	"chatgpt-discord-bot/src/config"
	"errors"
	"fmt"
	"log"
	"sort"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
)

type StartChatInput struct {
	Prompt    string
	Session   *discordgo.Session
	ChannelID string
	Config    *config.Config
	UserID    string
}

type ReplyChatInput struct {
	Prompt    string
	Session   *discordgo.Session
	ChannelID string
	Config    *config.Config
}

const GENERIC_REPLY string = "Chat started, reply in thread %s"
const CHATGPT_MODEL string = "gpt-3.5-turbo"
const BOT_USER_ID string = "1068334708057981022"

func StartChat(input StartChatInput) error {
	ctx := context.Background()

	log.Println("Starting thread")
	initialPrompt := input.Prompt
	thread, err := input.Session.ForumThreadStart(input.ChannelID, initialPrompt, 60, initialPrompt)

	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Println("Adding user to thread")
	input.Session.ThreadMemberAdd(thread.ID, input.UserID)

	session := openai.NewSession(input.Config.OpenAIApiKey)
	client := chat.NewClient(session, CHATGPT_MODEL)

	params := chat.CreateCompletionParams{
		Messages: []*chat.Message{
			{
				Role:    "user",
				Content: input.Prompt,
			},
		},
		N:         1,
		MaxTokens: 200,
	}

	log.Println("Sending request to OpenAI")
	resp, err := client.CreateCompletion(ctx, &params)

	if err != nil {
		fmt.Println(err)
		return err
	}

	message := resp.Choices[0].Message.Content

	log.Println("Sending OpenAI Response")
	input.Session.ChannelMessageSend(thread.ID, message)

	return nil
}

func ReplyInChat(input ReplyChatInput) error {
	var thread *discordgo.Channel
	var messages []*chat.Message = make([]*chat.Message, 0)

	ctx := context.Background()

	thread, err := input.Session.Channel(input.ChannelID)

	if err != nil {
		log.Fatalf("Error fetching threads: %s", err.Error())
		return err
	}

	if !thread.IsThread() {
		log.Fatalf("Channel %s is not a thread", input.ChannelID)

		return errors.New("Channel is not a thread. Please use this command in a thread started by ChatGPT")
	}

	threadMessages, err := input.Session.ChannelMessages(input.ChannelID, 60, "", "", "")

	if err != nil {
		log.Fatalf("Error fetching messages: %s", err.Error())
		return err
	}

	sort.Slice(threadMessages, func(i, j int) bool {
		return threadMessages[i].Timestamp.Before(threadMessages[j].Timestamp)
	})

	for _, message := range threadMessages {
		var role string = "user"
		fmt.Println(message.Content)
		if message.Content == "" || message.Content == "Generating reply" {
			continue
		}

		if message.Author.ID == BOT_USER_ID {
			role = "assistant"
		}

		messages = append(messages, &chat.Message{
			Role:    role,
			Content: message.Content,
		})
	}

	// append user's reply
	messages = append(messages, &chat.Message{
		Role:    "user",
		Content: input.Prompt,
	})

	session := openai.NewSession(input.Config.OpenAIApiKey)
	client := chat.NewClient(session, CHATGPT_MODEL)

	params := chat.CreateCompletionParams{
		Messages:  messages,
		N:         1,
		MaxTokens: 200,
	}

	resp, err := client.CreateCompletion(ctx, &params)

	if err != nil {
		fmt.Println(err)
		return err
	}

	message := resp.Choices[0].Message.Content
	input.Session.ChannelMessageSend(input.ChannelID, message)

	return nil
}
