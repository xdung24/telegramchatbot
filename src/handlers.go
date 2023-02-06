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
	gpt := GPTRepository{}
	gcloud := GoogleCloud{}

	question := c.Text()
	fmt.Println("Q:" + question)
	resp, err := gpt.GetGPTTextAnswer(question)
	if err != nil {
		return err
	}
	answer := resp.Choices[0].Text
	// trim space and new line characters
	answer = strings.TrimSpace(strings.TrimSuffix(answer, "\n"))

	// Will send awnser as message
	fmt.Println("A:" + answer)
	c.Send(answer)

	// Will send answer as audio file
	d, e := gcloud.DetectLanguage(question)
	if e != nil {
		fmt.Println(e.Error())
	}
	lang := "vi"
	if d != nil && d.Language.String() != "und" {
		fmt.Println("lang:" + d.Language.String())
		lang = d.Language.String()
	}

	file := gcloud.Prompt2Audio(answer, lang)
	audio := &tele.Audio{File: tele.FromDisk(file)}
	return c.Send(audio)
}
