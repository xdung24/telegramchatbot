package src

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const OnStartMessage = "Ask me something."
const MessagePrefix = "ChatGPT"

type Config struct {
	GPTModel          string
	GPTTemperature    float32
	GPTMaxTokens      int
	GPTCompletionsUrl string
}

var DefaultConfig = &Config{
	GPTModel:          "text-davinci-003",
	GPTTemperature:    0.5,
	GPTMaxTokens:      256,
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
	// TODO: Replace any for a custom interface
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
