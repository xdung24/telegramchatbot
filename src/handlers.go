package src

import (
	tele "gopkg.in/telebot.v3"
)

func OnStart(c tele.Context) error {
	return c.Send(OnStartMessage)
}

type GPTHandle struct{}

func (h *GPTHandle) AskGPT(c tele.Context) error {
	repo := GPTRepository{}
	resp, err := repo.GetGPTTextAnswer(c.Text())
	if err != nil {
		return err
	}
	return c.Send(resp.Choices[0].Text)
}
