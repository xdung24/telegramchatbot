package src

type GPTRepository struct{}

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
