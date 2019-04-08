package flipapps

import (
	fmt "fmt"
	"image/color"
	"image/draw"

	"image"
	"time"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/briggySmalls/flipdot/app/text"
)

type Imager interface {
	Message(sender, message string) ([]*flipdot.Image, error)
	Clock(time time.Time, isMessagesAvailable bool) ([]*flipdot.Image, error)
}

type imager struct {
	builder     text.TextBuilder
	signCount   uint
	statusImage image.Image
}

func NewImager(builder text.TextBuilder, statusImage image.Image, signCount uint) Imager {
	return &imager{
		builder:     builder,
		signCount:   signCount,
		statusImage: statusImage,
	}
}

// Helper function to send text to the signs
func (i *imager) Message(sender, message string) (images []*flipdot.Image, err error) {
	// Convert the sender to images
	senderImages, err := i.builder.Images(fmt.Sprintf("From: %s", sender), true)
	if err != nil {
		return
	}
	// Add empty images to fill frame, if necessary
	for uint(len(senderImages))%i.signCount != 0 {
		var emptyImage []draw.Image
		emptyImage, err = i.builder.Images("", false)
		if err != nil {
			return
		}
		senderImages = append(senderImages, emptyImage[0])
	}
	// Convert the text to images
	messageImages, err := i.builder.Images(message, true)
	if err != nil {
		return
	}
	// Combine the image sets
	srcImages := append(senderImages, messageImages...)
	// Convert the images to C form
	images = convertImages(srcImages)
	return
}

func (i *imager) Clock(time time.Time, isMessagesAvailable bool) (images []*flipdot.Image, err error) {
	// Get images that represent the time
	srcImages, err := i.builder.Images(time.Format("Mon 1 Jan\n3:04:05 pm"), true)
	errorHandler(err)
	// Add status if necessary
	if isMessagesAvailable {
		// Get far-right area the size of status image
		xOffset := srcImages[0].Bounds().Dx() - i.statusImage.Bounds().Dx()
		drawBounds := i.statusImage.Bounds().Add(image.Point{X: xOffset, Y: 0})
		draw.Draw(srcImages[0], drawBounds, i.statusImage, image.Point{}, draw.Over)
	}
	// Convert and return images
	images = convertImages(srcImages)
	return
}

// Packs an image into a C-style boolean array
func Slice(image image.Image) []bool {
	bgColor := color.Gray{0}
	// Create an array for the image
	rows := image.Bounds().Dy()
	cols := image.Bounds().Dx()
	binImage := make([]bool, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			binImage[r*cols+c] = image.At(c, r) != bgColor
		}
	}
	return binImage
}

func convertImages(inImages []draw.Image) (outImages []*flipdot.Image) {
	for _, img := range inImages {
		outImages = append(outImages, &flipdot.Image{Data: Slice(img)})
	}
	return
}
