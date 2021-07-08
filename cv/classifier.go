package cv

import (
	"errors"
	"gocv.io/x/gocv"
	"log"
	"os"
)

func GetClassifier() (gocv.CascadeClassifier, error) {
	classifier := gocv.NewCascadeClassifier()

	cascadeFileLink, exists := os.LookupEnv("CASCADE_FILE")
	if !exists {
		log.Panic("Can't start bot without cascade file")
	}

	if !classifier.Load(cascadeFileLink) {
		return classifier, errors.New("cascade file couldn't be loaded")
	}

	return classifier, nil
}
