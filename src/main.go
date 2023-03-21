package main

import (
	"chatgpt-discord-bot/src/chat"
	config "chatgpt-discord-bot/src/config"
	verify "chatgpt-discord-bot/src/verify"
	"errors"
	"fmt"

	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session
var appConfig *config.Config

func init() {
	var err error

	appConfig = config.CreateConfig()
	session, err = discordgo.New(fmt.Sprintf("Bot %s", appConfig.DiscordToken))

	if err != nil {
		log.Fatal(err.Error())
	}
}

type DiscordInteractionRequest struct {
	Headers DiscordInteractionHeaders `json:"headers"`
	Body    string                    `json:"body"`
}

type DiscordInteractionHeaders struct {
	Signature string `json:"X-Signature-Ed25519"`
	Timestamp string `json:"X-Signature-Timestamp"`
}

func HandleRequest(ctx context.Context, request DiscordInteractionRequest) (discordgo.InteractionResponse, error) {
	var interaction discordgo.Interaction

	fmt.Println(request)

	verifyInput := verify.VerifyInput{
		Signature: request.Headers.Signature,
		Body:      request.Body,
		Timestamp: request.Headers.Timestamp,
		Config:    appConfig,
	}

	if !verify.Verify(verifyInput) {
		log.Fatal("Could not verify signature")

		return discordgo.InteractionResponse{}, errors.New("Could not verify signature")
	}

	log.Println("unmarshalling body")
	err := interaction.UnmarshalJSON([]byte(request.Body))

	if err != nil {
		log.Fatalf("Could not decode body: %s\n", err.Error())

		return discordgo.InteractionResponse{}, err

	}

	if interaction.Type == discordgo.InteractionPing {
		return discordgo.InteractionResponse{
			Type: 1,
		}, nil
	}
	var command discordgo.ApplicationCommandInteractionData = interaction.ApplicationCommandData()

	if err != nil {
		log.Fatalf("Could not decode body: %s\n", err.Error())

		return discordgo.InteractionResponse{}, err

	}

	switch command.Name {
	case "start-chat":
		session.InteractionRespond(&interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Generating chat",
				TTS:     false,
			},
		})

		input := chat.StartChatInput{
			Prompt:    command.Options[0].StringValue(),
			ChannelID: interaction.ChannelID,
			UserID:    interaction.Member.User.ID,
			Session:   session,
			Config:    appConfig,
		}

		err = chat.StartChat(input)

		if err != nil {
			log.Fatalf(err.Error())

			return discordgo.InteractionResponse{}, err
		}

		log.Println("Returning response")
		return discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Response generated, see created thread",
				TTS:     false,
			},
		}, nil
	case "reply":
		session.InteractionRespond(&interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Generating reply",
				TTS:     false,
			},
		})
		input := chat.ReplyChatInput{
			Prompt:    command.Options[0].StringValue(),
			ChannelID: interaction.ChannelID,
			Session:   session,
			Config:    appConfig,
		}

		err = chat.ReplyInChat(input)

		if err != nil {
			log.Fatalf(err.Error())

			return discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
					TTS:     false,
				},
			}, err
		}

		log.Println("Returning response")
		return discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Reply generated",
				TTS:     false,
			},
		}, nil
	}

	return discordgo.InteractionResponse{}, errors.New("Interaction not found")
}

func main() {
	lambda.Start(HandleRequest)
}
