package chathandler

import (
	"chatgpt-discord-bot/src/config"
	"fmt"

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
