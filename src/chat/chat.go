package chathandler

import (
	"chatgpt-discord-bot/src/config"
	"fmt"
	"sort"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
)

type CommandHandlerInput struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	Config      *config.Config
}

type ReplyChatInput struct {
	Session *discordgo.Session
	Thread  *discordgo.Channel
	Config  *config.Config
}

const GENERIC_REPLY string = "Starting Chat..."
const CHATGPT_MODEL string = "gpt-3.5-turbo"

func StartChat(input CommandHandlerInput) error {
	ctx := context.Background()

	initialPrompt := input.Interaction.ApplicationCommandData().Options[0].StringValue()
	thread, err := input.Session.ForumThreadStart(input.Interaction.ChannelID, initialPrompt, 60, initialPrompt)

	if err != nil {
		fmt.Println(err)
		return err
	}

	input.Session.ThreadMemberAdd(thread.ID, input.Interaction.Member.User.ID)

	session := openai.NewSession(input.Config.OpenAIApiKey)
	client := chat.NewClient(session, CHATGPT_MODEL)

	params := chat.CreateCompletionParams{
		Messages: []*chat.Message{
			{
				Role:    "user",
				Content: input.Interaction.ApplicationCommandData().Options[0].StringValue(),
			},
		},
		N:         1,
		MaxTokens: 200,
	}

	resp, err := client.CreateCompletion(ctx, &params)

	if err != nil {
		fmt.Println(err)
		return err
	}

	message := resp.Choices[0].Message.Content
	fmt.Println(thread.Messages)
	input.Session.ChannelMessageSend(thread.ID, message)

	return nil
}

func ReplyInChat(input ReplyChatInput) error {
	var messages []*chat.Message = make([]*chat.Message, 0)

	ctx := context.Background()

	threadMessages, err := input.Session.ChannelMessages(input.Thread.ID, 100, "", "", "")

	sort.Slice(threadMessages, func(i, j int) bool {
		return threadMessages[i].ID < threadMessages[j].ID
	})
	for _, message := range threadMessages {
		fmt.Println(message.Content)
		var role string = "user"

		if message.Content == "" {
			continue
		}

		if message.Author.ID == input.Session.State.User.ID {
			role = "assistant"
		}

		messages = append(messages, &chat.Message{
			Role:    role,
			Content: message.Content,
		})
	}

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
	input.Session.ChannelMessageSend(input.Thread.ID, message)

	return nil
}
