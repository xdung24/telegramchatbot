package src

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

type GPTRepository struct {
	OpenApiKey string
	Config     Config
}

const gptCompletionsUrl = "https://api.openai.com/v1/completions"

func (r *GPTRepository) GetGPTTextAnswer(prompt string) (*TextResponse, error) {
	textRequest := TextRequest{
		Model:       r.Config.GPTModel,
		Prompt:      prompt,
		Temperature: r.Config.GPTTemperature,
		MaxTokens:   r.Config.GPTMaxTokens,
	}

	res, err := call("POST", gptCompletionsUrl, textRequest, r.OpenApiKey)
	if err != nil {
		return nil, err
	}

	textResponse := TextResponse{}
	err = readResponse(res, &textResponse)
	if err != nil {
		return nil, err
	}

	return &textResponse, nil
}

func call(method string, url string, body any, openApiKey string) (*http.Response, error) {
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readResponse(res *http.Response, model any) error {
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
