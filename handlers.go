package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosimple/slug"
	"github.com/xdung24/gpt-bot/src"
	tele "gopkg.in/telebot.v3"
)

const OnStartMessage = "Ask me something. Please prefer English, but other languages are also supported."

type Handler struct {
	GPTRepository   src.GPTRepository
	GoogleCloud     src.GoogleCloud
	ParagraphLength int
}

func (h *Handler) OnStart(c tele.Context) error {
	return c.Send(OnStartMessage)
}

func (h *Handler) AskGPT(c tele.Context) error {

	// print input question
	question := c.Text()
	fmt.Println("Q1:" + question)

	// detect question language
	d, e := h.GoogleCloud.DetectLanguage(question)
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
		result, err := h.GoogleCloud.TranslateText("en", question)
		if err != nil {
			fmt.Println("error:" + err.Error())
		}
		question = result
		fmt.Println("Q2:" + question)
	}

	// get chat gpt completion (answer)
	resp, err := h.GPTRepository.GetGPTTextAnswer(question)
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

	// split answer to paragraphs
	var answerList []string
	if len(answer) > h.ParagraphLength {
		answers := src.SplitParagraphs(answer, h.ParagraphLength)
		answerList = append(answerList, answers...)
	} else {
		answerList = append(answerList, answer)
	}

	// translate answer to original language
	for _, answer := range answerList {

		if questionlang != "en" {
			result, err := h.GoogleCloud.TranslateText(questionlang, answer)
			if err != nil {
				fmt.Println("error:" + err.Error())
			}
			answer = result
			fmt.Println("A2:" + answer)
		}

		file, err := h.GoogleCloud.Prompt2Audio(answer, questionlang)
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
		c.Send(audio)
	}
	return nil
}
