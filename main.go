package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/johnazedo/gpt-bot/src"
	"github.com/joho/godotenv"
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
	handle := src.GPTHandle{}

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_KEY"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", src.OnStart)
	b.Handle("/version", Version)
	b.Handle(tele.OnText, handle.AskGPT)

	b.Start()
}
