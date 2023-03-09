package main

import (
	"fmt"
	"html"
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
		sendErr := c.Send(err.Error())
		logSendErr(sendErr)
		return sendErr
	}
	if len(resp.Choices) == 0 {
		sendErr := c.Send("No answer")
		logSendErr(sendErr)
		return sendErr
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
	for i, answer := range answerList {

		var unescapedAnswer string
		if questionlang != "en" {
			result, err := h.GoogleCloud.TranslateText(questionlang, answer)
			if err != nil {
				fmt.Println("error:" + err.Error())
			}
			unescapedAnswer = html.UnescapeString(result)
			fmt.Printf("A2(%d): %s\n", i, unescapedAnswer)
		} else {
			unescapedAnswer = html.UnescapeString(answer)
		}

		file, err := h.GoogleCloud.Prompt2Audio(unescapedAnswer, questionlang)
		if err != nil {
			// send answer as message
			sendErr := c.Send(unescapedAnswer)
			logSendErr(sendErr)
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
		sendErr := c.Send(audio)
		logSendErr(sendErr)
	}
	return nil
}

func logSendErr(err error) {
	if err != nil {
		fmt.Println("error when sending message" + err.Error())
	} else {
		fmt.Println("send message successfully")
	}
}
