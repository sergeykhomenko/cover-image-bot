package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	tgApiKey, exists := os.LookupEnv("TG_API_KEY")
	if !exists {
		log.Panic("Can't start bot without API key")
	}

	bot, err := tgbotapi.NewBotAPI(tgApiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		replyTo := update.Message.ReplyToMessage
		if replyTo == nil || replyTo.From.IsBot == false {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
