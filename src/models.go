package src

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
