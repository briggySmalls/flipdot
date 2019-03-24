package flipapps

import (
	fmt "fmt"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/briggySmalls/flipdot/app/text"
)

type Imager interface {
	Message(sender, message string) ([]*flipdot.Image, error)
	Clock() ([]*flipdot.Image, error)
}

type imager struct {
	builder   text.TextBuilder
	signCount uint
}

func NewImager(builder text.TextBuilder, statusImage Image, signCount uint) {
	return &imager{
		builder:   builder,
		signCount: signCount,
		statusImage: Image,
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
	for len(senderImages)%i.signCount != 0 {
		senderImages = append(senderImages, i.builder.Images("", false))
	}
	// Convert the text to images
	messageImages, err := textBuilder.Images(txt, true)
	if err != nil {
		return
	}
	// Combine the image sets
	images := append(senderImages, messageImages...)
	// Convert the images to C form
	var packedImages []*flipdot.Image
	for _, img := range images {
		packedImages = append(packedImages, &flipdot.Image{Data: Slice(img)})
	}
	// Send the text
	err = s.flipdot.Draw(packedImages)
	return
}

func (i *imager) Clock(time time.Time, isMessagesAvailable bool) {
	// Get images that represent the time
	images, err := i.builder.Images(t.Format("Mon 1 Jan\n3:04 pm"), true)
	// Add status if necessary
	if isMessagesAvailable {
		images[0]
	}
}

func readStatusImage(filename string) (image Image, err error) {
	// Read the image from disk
	file, err := os.Open(filename)
    if err != nil {
        return
    }
	// Interpret as an image
	image, err := png.Decode(file)
    if err != nil {
        return
    }
	return
}
