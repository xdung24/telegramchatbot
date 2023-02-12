package main

import (
	"fmt"
	"log"
	"os"
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
	fmt.Println("version=", version)
	fmt.Println("build=", build)

	_ = godotenv.Load()
	gpt := src.GPTHandle{}

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_KEY"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", src.OnStart)
	bot.Handle("/version", Version)
	bot.Handle(tele.OnText, gpt.AskGPT)

	bot.Start()
}
