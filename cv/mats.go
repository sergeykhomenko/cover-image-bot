package cv

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type CVImagePrepared struct {
	mat    gocv.Mat
	Width  int
	Height int
	Faces  []image.Rectangle
}

func NewImagePrepared(filename string) CVImagePrepared {
	var img CVImagePrepared

	file, _ := os.Open(filename)
	defer file.Close()

	decodedFileImage, _, _ := image.Decode(file)
	img.convertImage(decodedFileImage)

	return img
}

func (cvImage *CVImagePrepared) convertImage(img image.Image) {
	bounds := img.Bounds()
	cvImage.Width = bounds.Dx()
	cvImage.Height = bounds.Dy()

	bytes := make([]byte, 0, cvImage.Width*cvImage.Height)
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8))
			bytes = append(bytes, byte(g>>8))
			bytes = append(bytes, byte(r>>8))
		}
	}

	cvImage.mat, _ = gocv.NewMatFromBytes(cvImage.Height, cvImage.Width, gocv.MatTypeCV8UC3, bytes)
}

func (cvImage *CVImagePrepared) SavePreparedImageToFile(filename string) {
	blue := color.RGBA{0, 0, 255, 0}

	// draw a rectangle around each face on the original image
	for _, r := range cvImage.Faces {
		gocv.Rectangle(&cvImage.mat, r, blue, 1)
	}

	gocv.IMWrite(filename+"-result.jpg", cvImage.mat)

	fmt.Println(filename + " - done")
}

func (cvImage *CVImagePrepared) DetectFaces(classifier gocv.CascadeClassifier) {
	cvImage.Faces = classifier.DetectMultiScale(cvImage.mat)
}
