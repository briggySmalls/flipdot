package client

import (
	context "context"
	fmt "fmt"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/text"
)

const (
	contextTimeoutS = 10
)

type Flipdot interface {
	Signs() []*GetInfoResponse_SignInfo
	Size() (width, height uint)
	LightOn() error
	LightOff() error
	TestStart() error
	TestStop() error
	Draw(images []*Image, isWait bool) error
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
	// Duration to space out message frames
	frameTime time.Duration
}

func NewFlipdot(client FlipdotClient, frameTime time.Duration) (f Flipdot, err error) {
	flipdot := flipdot{
		client:    client,
		frameTime: frameTime,
	}
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

// Get info from the sign
func (f *flipdot) Signs() (signs []*GetInfoResponse_SignInfo) {
	return f.signs
}

// Return accepted size for flipdot
func (f *flipdot) Size() (width, height uint) {
	// Get the first sign
	sign := f.Signs()[0]
	// Return the sign's dimensions
	return uint(sign.Width), uint(sign.Height)
}

// Draw a set of images
func (f *flipdot) Draw(images []*Image, isWait bool) (err error) {
	// Send any relevant images
	images, err = f.sendFrame(images)
	if err != nil || (len(images) == 0 && !isWait) {
		// Either:
		// 1) We've errored
		// 2) We've sent all our images and we've not been asked to keep them visible
		return
	}
	// Create a ticker
	ticker := time.NewTicker(f.frameTime)
	defer ticker.Stop()
	// Write images periodically
	for {
		select {
		case <-ticker.C:
			if len(images) > 0 {
				// Send a frame's-worth of images
				images, err = f.sendFrame(images)
				if err != nil {
					return
				}
			} else {
				// We've finished displaying images, move on
				return
			}
		}
	}
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
	_, err = f.client.Light(ctx, &LightRequest{Status: status})
	// Handle errors
	return
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
	_, err = f.client.Test(ctx, &TestRequest{Action: action})
	return
}

// Send a set of images to available signs
func (f *flipdot) sendFrame(images []*Image) (leftover []*Image, err error) {
	leftover = images
	width, height := f.Size()
	blankImageData := make([]bool, width*height)
	for _, sign := range f.signNames {
		// Send an empty image if there are none left (removes old messages)
		if len(leftover) == 0 {
			f.writeImage(Image{Data: blankImageData}, sign)
			return
		}
		// Pop an image off the stack and send it
		var image *Image
		image, leftover = leftover[0], leftover[1:]
		err = f.writeImage(*image, sign)
		if err != nil {
			return
		}
	}
	return leftover, nil
}

// Write an image to the specified sign
func (f *flipdot) writeImage(image Image, sign string) (err error) {
	// Send request
	ctx, cancel := getContext()
	defer cancel()
	_, err = f.client.Draw(ctx, &DrawRequest{
		Sign:  sign,
		Image: &image,
	})
	return
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

// Request signs information from service
func (f *flipdot) getSigns() (signs []*GetInfoResponse_SignInfo, err error) {
	// Get the signs
	context, cancel := getContext()
	defer cancel()
	response, err := f.client.GetInfo(context, &GetInfoRequest{})
	if err != nil {
		// Something went wrong
		return nil, err
	}
	return response.Signs, nil
}

// Get a simple context for sending requests via gRPC
func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeoutS*time.Second)
}
