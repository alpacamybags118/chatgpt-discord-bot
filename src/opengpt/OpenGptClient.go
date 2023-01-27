package opengptclient

import (
	"bytes"
	"chatgpt-discord-bot/src/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenGptClient struct {
	config config.Config
}

type OpenGptCompletionRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
	Max_tokens  int     `json:"max_tokens"`
}

type OpenGptCompletionResponse struct {
	Choices []OpenGptCompletionChoice
}

type OpenGptCompletionChoice struct {
	Text          string
	Index         int
	Finish_reason string
}

func CreateNew(config *config.Config) OpenGptClient {
	return OpenGptClient{
		config: *config,
	}
}

func (c OpenGptClient) SendCompletionRequest(request OpenGptCompletionRequest) (OpenGptCompletionResponse, error) {
	var completion OpenGptCompletionResponse
	client := &http.Client{}

	reqBody, err := json.Marshal(request)

	if err != nil {
		return completion, err
	}
	fmt.Println(string(reqBody))
	body := bytes.NewReader(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/completions", c.config.ChatGPTUrl), body)

	if err != nil {
		return completion, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.OpenAIApiKey))

	resp, err := client.Do(req)

	if err != nil {
		return completion, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return completion, err
	}

	err = json.Unmarshal(respBody, &completion)
	fmt.Println(resp.Status)
	if err != nil {
		return completion, err
	}

	return completion, nil
}
