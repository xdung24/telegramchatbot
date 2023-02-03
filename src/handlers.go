package src

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func OnStart(c tele.Context) error {
	return c.Send(OnStartMessage)
}

type GPTHandle struct{}

func (h *GPTHandle) AskGPT(c tele.Context) error {
	repo := GPTRepository{}
	question := c.Text()
	fmt.Println("Q:" + question)
	resp, err := repo.GetGPTTextAnswer(question)
	if err != nil {
		return err
	}
	answer := resp.Choices[0].Text
	// trim space and new line characters
	answer = strings.TrimSpace(strings.TrimSuffix(answer, "\n"))
	fmt.Println("A:" + answer)
	return c.Send(answer)
}
