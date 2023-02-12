package src

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosimple/slug"
	tele "gopkg.in/telebot.v3"
)

func OnStart(c tele.Context) error {
	return c.Send(OnStartMessage)
}

type GPTHandle struct{}

func (h *GPTHandle) AskGPT(c tele.Context) error {
	gpt := GPTRepository{}
	gcloud := GoogleCloud{}

	// print input question
	question := c.Text()
	fmt.Println("Q1:" + question)

	// detect question language
	d, e := gcloud.DetectLanguage(question)
	if e != nil {
		fmt.Println(e.Error())
	}
	questionlang := ""
	if d != nil && d.Language.String() != "und" {
		fmt.Println("lang:" + d.Language.String())
		questionlang = d.Language.String()
	}

	// not english, translate question to english
	if questionlang != "en" {
		result, err := gcloud.TranslateText("en", question)
		if err != nil {
			fmt.Println("error:" + err.Error())
		}
		question = result
		fmt.Println("Q2:" + question)
	}

	// get chat gpt completion (answer)
	resp, err := gpt.GetGPTTextAnswer(question)
	if err != nil {
		return c.Send(err.Error())
	}
	if len(resp.Choices) == 0 {
		return c.Send("No answer")
	}

	answer := resp.Choices[0].Text

	// trim space and new line characters
	answer = strings.TrimSpace(answer)
	answer = strings.ReplaceAll(answer, "\n\n", "\n")
	answer = strings.ReplaceAll(answer, "\n", ".\n")
	answer = strings.ReplaceAll(answer, "..\n", ".\n")

	fmt.Println("A1:" + answer)

	// translate answer to original language
	if questionlang != "en" {
		result, err := gcloud.TranslateText(questionlang, answer)
		if err != nil {
			fmt.Println("error:" + err.Error())
		}
		answer = result
		fmt.Println("A2:" + answer)
	}

	file, err := gcloud.Prompt2Audio(answer, questionlang)
	if err != nil {
		// send answer as message
		return c.Send(answer)
	}

	filename := slug.Make(question) + ".ogg"

	audio := &tele.Audio{
		File:     tele.FromDisk(file),
		FileName: filename,
		Title:    question,
		Caption:  answer,
	}
	defer os.Remove(file)
	// send answer as voice message
	return c.Send(audio)
}
