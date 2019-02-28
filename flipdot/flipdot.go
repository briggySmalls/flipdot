package flipdot

import (
	context "context"
	fmt "fmt"
	"image"
	"time"

	"github.com/briggySmalls/flipcli/text"
	"golang.org/x/image/font"
)

const (
	contextTimeoutS = 10
)

type Flipdot interface {
	LightOn() error
	LightOff() error
	TestStart() error
	TestStop() error
	Draw(images []image.Image) error
	Text(text string, font font.Face) error
}

type flipdot struct {
	// gRPC client to send command via
	client FlipdotClient
	// Record of response from GetInfo request
	signs []*GetInfoResponse_SignInfo
	// Names of signs from GetInfo request
	signNames []string
	// TextBuilder used to convert text to images
	textBuilder text.TextBuilder
}

func NewFlipdot(client FlipdotClient) (f Flipdot, err error) {
	flipdot := flipdot{client: client}
	err = flipdot.init()
	f = Flipdot(&flipdot)
	return
}

// Send request to turn the light on
func (f *flipdot) LightOn() (err error) {
	return f.light(true)
}

// Send request to turn the light off
func (f *flipdot) LightOff() (err error) {
	return f.light(false)
}

// Send request to start the test sequence
func (f *flipdot) TestStart() (err error) {
	return f.test(true)
}

// Send request to turn the light off
func (f *flipdot) TestStop() (err error) {
	return f.test(false)
}

// Draw a set of images
func (f *flipdot) Draw(images []image.Image) (err error) {
	// Convert stdlib images to our type
	// TODO: Move our image interface here and do away with img type
	var imgs []text.Image
	for _, image := range images {
		imgs = append(imgs, text.NewImage(image))
	}
	// Send the images
	err = f.sendImages(imgs)
	return
}

// Draw some text
func (f *flipdot) Text(txt string, font font.Face) (err error) {
	// Create a text builder
	f.textBuilder = text.NewTextBuilder(uint(f.signs[0].Width), uint(f.signs[0].Height), font)
	// Convert the text to images
	images, err := f.textBuilder.Images(txt)
	if err != nil {
		return
	}
	// Send the text
	err = f.sendImages(images)
	return
}

// Initialise the struct with some one-off attributes
func (f *flipdot) init() (err error) {
	// Get the signs for later
	f.signs, err = f.getSigns()
	if err != nil {
		return
	}
	// Validate the signs
	err = checkSigns(f.signs)
	if err != nil {
		return
	}
	// Get the sign names
	var signNames []string
	for _, sign := range f.signs {
		signNames = append(signNames, sign.Name)
	}
	f.signNames = signNames
	return
}

// Send request to set the light status
func (f *flipdot) light(on bool) (err error) {
	// Get context
	ctx, cancel := getContext()
	defer cancel()
	// Send request
	var status LightRequest_Status
	if on {
		status = LightRequest_ON
	} else {
		status = LightRequest_OFF
	}
	response, err := f.client.Light(ctx, &LightRequest{Status: status})
	// Handle errors
	return handleErrors(err, response.Error)
}

// Send request to start/stop test sequence
func (f *flipdot) test(start bool) (err error) {
	// Get context
	ctx, cancel := getContext()
	defer cancel()
	// Send request
	var action TestRequest_Action
	if start {
		action = TestRequest_START
	} else {
		action = TestRequest_STOP
	}
	response, err := f.client.Test(ctx, &TestRequest{Action: action})
	return handleErrors(err, response.Error)
}

// Send a set of images, periodically if necessary
func (f *flipdot) sendImages(images []text.Image) (err error) {
	// Send any relevant images
	images, err = f.sendFrame(images)
	// Check if we need to go on
	if len(images) == 0 || err != nil {
		return
	}
	// Create a ticker
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	// Write images periodically
	for len(images) > 0 {
		select {
		case <-ticker.C:
			images, err = f.sendFrame(images)
			if err != nil {
				return
			}
		}
	}
	return
}

// Send a set of images to available signs
func (f *flipdot) sendFrame(images []text.Image) (leftover []text.Image, err error) {
	for _, sign := range f.signNames {
		// Stop sending if there are no more images left
		if len(images) == 0 {
			return
		}
		// Pop an image off the stack and send it
		var image text.Image
		image, images = images[0], images[1:]
		err = f.writeImage(image, sign)
		if err != nil {
			return
		}
	}
	return
}

// Write an image to the specified sign
func (f *flipdot) writeImage(image text.Image, sign string) (err error) {
	// Send request
	ctx, cancel := getContext()
	defer cancel()
	response, clientError := f.client.Draw(ctx, &DrawRequest{
		Sign:  sign,
		Image: image.Slice(),
	})
	return handleErrors(clientError, response.Error)
}

// Check that all signs have the same width/height
func checkSigns(signs []*GetInfoResponse_SignInfo) error {
	var width, height uint32
	for i, sign := range signs {
		if i == 0 {
			width = sign.Width
			height = sign.Height
		} else {
			if width != sign.Width {
				return fmt.Errorf("Sign width %d != %d", sign.Width, width)
			} else if height != sign.Height {
				return fmt.Errorf("Sign height %d != %d", sign.Height, height)
			}
		}
	}
	return nil
}

func (f *flipdot) getSigns() (signs []*GetInfoResponse_SignInfo, err error) {
	// Get the signs
	context, cancel := getContext()
	defer cancel()
	response, clientErr := f.client.GetInfo(context, &GetInfoRequest{})
	err = handleErrors(clientErr, response.Error)
	if err != nil {
		// Something went wrong
		return nil, err
	}
	return response.Signs, nil
}

// Handle error and response from gRPC call
func handleErrors(clientError error, serverError *Error) error {
	if clientError != nil {
		return clientError
	}
	if serverError.Code != 0 {
		return toError(serverError)
	}
	return nil
}

// Convert a server error message (sent via gRPC) to an error
func toError(serverError *Error) (err error) {
	return fmt.Errorf("Flipdot server error (%d): %s", serverError.Code, serverError.Message)
}

// Get a simple context for sending requests via gRPC
func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeoutS*time.Second)
}
