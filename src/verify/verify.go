package verify

import (
	"bytes"
	"chatgpt-discord-bot/src/config"
	"crypto/ed25519"
	"encoding/hex"
	"log"
)

type VerifyInput struct {
	Signature string
	Body      string
	Timestamp string
	Config    *config.Config
}

// Verify verifies the signature provided in the request
// with the public key of the bot
func Verify(input VerifyInput) bool {
	var body bytes.Buffer

	decodedKey, err := hex.DecodeString(input.Config.DiscordPublicKey)

	if err != nil {
		log.Fatalf("Error when decoding public key %s\n", err.Error())

		return false
	}

	publicKey := ed25519.PublicKey(decodedKey)

	signature, err := hex.DecodeString(input.Signature)

	if err != nil {
		log.Fatalf("Error when decoding signature %s\n", err.Error())

		return false
	}

	if len(signature) != ed25519.SignatureSize || signature[63]&224 != 0 {
		log.Fatal("Signature is invalid, signature size is incorrect")

		return false
	}

	body.WriteString(input.Timestamp)
	body.WriteString(input.Body)

	return ed25519.Verify(publicKey, body.Bytes(), signature)
}
