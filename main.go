package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/sergeykhomenko/cover-image-bot/cv"
	"gocv.io/x/gocv"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// load classifier to recognize faces
	classifier, err := cv.GetClassifier()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer classifier.Close()

	prepareDataset(classifier)
}

func prepareDataset(classifier gocv.CascadeClassifier) {
	// load image list to recognize
	images, _ := ioutil.ReadDir("./data/dataset")

	for _, fileInfo := range images {
		if fileInfo.IsDir() {
			continue
		}

		filename := "./data/dataset/" + fileInfo.Name()

		img := cv.NewImagePrepared(filename)
		img.DetectFaces(classifier)
		img.SavePreparedImageToFile(filename)
	}
}

func tgmain() {
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

		//bot.GetFileDirectURL(update.Message.Photo[0].FileID)

		log.Printf("[%s] %s", update.Message.Photo)
	}
}
