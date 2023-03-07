package src

import (
	"context"
	"fmt"
	"os"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type GoogleCloud struct {
	Gender string
}

func ssmlFromString(genderText string) texttospeechpb.SsmlVoiceGender {
	if genderText == "FEMALE" || genderText == "2" {
		return texttospeechpb.SsmlVoiceGender_FEMALE
	}
	if genderText == "MALE" || genderText == "1" {
		return texttospeechpb.SsmlVoiceGender_MALE
	}
	panic("invalid genderText")
}

func (gc *GoogleCloud) Prompt2Audio(prompt string, lang string) (string, error) {
	// Instantiates a client.
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer client.Close()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: prompt},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: lang,
			SsmlGender:   ssmlFromString(gc.Gender),
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// The resp's AudioContent is binary.
	filename, _ := os.CreateTemp(os.TempDir(), "*.ogg")
	err = os.WriteFile(filename.Name(), resp.AudioContent, 0644)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return filename.Name(), nil
}

func (gc *GoogleCloud) DetectLanguage(text string) (*translate.Detection, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("translate.NewClient: %v", err)
	}
	defer client.Close()
	lang, err := client.DetectLanguage(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("detectLanguage: %v", err)
	}
	if len(lang) == 0 || len(lang[0]) == 0 {
		return nil, fmt.Errorf("detectLanguage return value empty")
	}
	return &lang[0][0], nil
}

func (gc *GoogleCloud) TranslateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
