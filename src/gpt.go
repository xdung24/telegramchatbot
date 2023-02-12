package src

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Config struct {
	GPTModel          string
	GPTTemperature    float32
	GPTMaxTokens      int
	GPTCompletionsUrl string
}

type TextRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type TextResponse struct {
	Id      string                `json:"id"`
	Created int                   `json:"created"`
	Choices []TextResponseChoices `json:"choices"`
}

type TextResponseChoices struct {
	Text     string `json:"text"`
	Index    int    `json:"index"`
	LogProbs int    `json:"logprobs"`
}

type GPTRepository struct{}

const OnStartMessage = "Ask me something."
const MessagePrefix = "ChatGPT"

var DefaultConfig = &Config{
	GPTModel:          "text-davinci-003",
	GPTTemperature:    0.7,
	GPTMaxTokens:      1000,
	GPTCompletionsUrl: "https://api.openai.com/v1/completions",
}

func Call(method, url string, body any) (*http.Response, error) {
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ReadResponse(res *http.Response, model any) error {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, model)
	if err != nil {
		return err
	}

	return nil
}

func (r *GPTRepository) GetGPTTextAnswer(prompt string) (*TextResponse, error) {
	textRequest := TextRequest{
		Model:       DefaultConfig.GPTModel,
		Prompt:      prompt,
		Temperature: DefaultConfig.GPTTemperature,
		MaxTokens:   DefaultConfig.GPTMaxTokens,
	}

	res, err := Call("POST", DefaultConfig.GPTCompletionsUrl, textRequest)
	if err != nil {
		return nil, err
	}

	textResponse := TextResponse{}
	err = ReadResponse(res, &textResponse)
	if err != nil {
		return nil, err
	}

	return &textResponse, nil
}
