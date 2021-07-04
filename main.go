package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//images := []string{
	//"./data/hq720.jpg",
	//"./data/photo_2021-07-03 214523.jpg",
	//"./data/photo_2021-07-03 21.45.29.jpg",
	//"./data/photo_2021-07-03 21.45.56.jpg",
	//}

	filename := "./data/photo_2021-07-03 21.45.29.jpg"

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	cascadeFileLink, exists := os.LookupEnv("CASCADE_FILE")
	if !exists {
		log.Panic("Can't start bot without cascade file")
	}

	if !classifier.Load(cascadeFileLink) {
		fmt.Println("Cascade file couldn't be loaded")
		return
	}

	window := gocv.NewWindow("test")

	mat := getImageMat(filename)
	defer mat.Close()

	blue := color.RGBA{0, 0, 255, 0}

	rects := classifier.DetectMultiScale(mat)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image
	for _, r := range rects {
		gocv.Rectangle(&mat, r, blue, 3)
	}

	for {
		window.IMShow(mat)
		window.WaitKey(1)
	}
}

func getImageMat(filename string) gocv.Mat {
	f, _ := os.Open(filename)
	defer f.Close()

	img, _, _ := image.Decode(f)
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()

	bytes := make([]byte, 0, x*y)
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8))
			bytes = append(bytes, byte(g>>8))
			bytes = append(bytes, byte(r>>8))
		}
	}

	mat, _ := gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)

	return mat
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
