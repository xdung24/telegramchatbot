package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/xdung24/gpt-bot/src"
	tele "gopkg.in/telebot.v3"
)

var (
	version string
	build   string
)

func Version(c tele.Context) error {
	return c.Send("version=" + version + "\n" + "build=" + build)
}

func main() {
	// Print version and build
	fmt.Println("version=", version)
	fmt.Println("build=", build)

	// Load environment variables
	_ = godotenv.Load()

	// Init GPT Repository
	var config = src.Config{
		GPTModel:       "text-davinci-003",
		GPTTemperature: 0.7,
		GPTMaxTokens:   4000,
	}
	gptRepo := src.GPTRepository{
		OpenApiKey: os.Getenv("OPENAI_KEY"),
		Config:     config,
	}

	// Init Google Cloud
	gCloud := src.GoogleCloud{
		Gender: os.Getenv("GENDER"),
	}

	// Init Handler with Google Cloud and GPT Repository and Paragraph Length
	paragraphLength, err := strconv.Atoi(os.Getenv("PARAGRAPH_LENGTH"))
	if err != nil {
		panic(err)
	}
	handler := Handler{
		GPTRepository:   gptRepo,
		GoogleCloud:     gCloud,
		ParagraphLength: paragraphLength,
	}

	// Init Telegram Bot setting
	settings := tele.Settings{
		Token:  os.Getenv("TELEGRAM_KEY"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	// Create Telegram Bot instance
	bot, err := tele.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle Telegram Bot commands
	bot.Handle("/version", Version)
	bot.Handle("/start", handler.OnStart)
	bot.Handle(tele.OnText, handler.AskGPT)

	// Start Telegram Bot
	bot.Start()
}
