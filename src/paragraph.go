package src

import (
	"strings"

	"github.com/neurosnap/sentences/english"
)

func SplitParagraphs(text string, maxChars int) []string {
	var paragraphs []string
	var currentParagraph strings.Builder
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	sentences := tokenizer.Tokenize(text)
	for _, sentence := range sentences {
		sentenceStr := sentence.Text
		sentenceLen := len(sentenceStr)
		if currentParagraph.Len() > 0 && currentParagraph.Len()+sentenceLen > maxChars {
			paragraphs = append(paragraphs, currentParagraph.String())
			currentParagraph.Reset()
		}
		if sentenceLen > maxChars {
			// Sentence is too long to fit in a single paragraph, split it up
			parts := splitString(sentenceStr, maxChars)
			for _, part := range parts {
				currentParagraph.WriteString(part)
				if currentParagraph.Len() >= maxChars {
					paragraphs = append(paragraphs, currentParagraph.String())
					currentParagraph.Reset()
				}
			}
		} else {
			currentParagraph.WriteString(sentenceStr)
		}
	}
	if currentParagraph.Len() > 0 {
		paragraphs = append(paragraphs, strings.TrimSpace(currentParagraph.String()))
	}
	return paragraphs
}

func splitString(str string, maxChars int) []string {
	var parts []string
	for len(str) > maxChars {
		index := strings.LastIndex(str[:maxChars+1], " ")
		if index < 0 {
			index = maxChars
		}
		parts = append(parts, str[:index])
		str = str[index+1:]
	}
	parts = append(parts, str)
	return parts
}
