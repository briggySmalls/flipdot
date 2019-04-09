package cmd

import (
	"context"
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/stianeikeland/go-rpio"
	"google.golang.org/grpc"
)

// Mock output pin
type mockOutputPin struct{}

func (p *mockOutputPin) High()   {}
func (p *mockOutputPin) Low()    {}
func (p *mockOutputPin) Toggle() {}

// Mock trigger pin
type mockTriggerPin struct {
	state bool // Button state
}

func (p *mockTriggerPin) Read() rpio.State {
	if p.state {
		return rpio.High
	} else {
		return rpio.Low
	}
}

// Mock flipdot
type mockFlipdotClient struct {
	signConfig   []*flipdot.GetInfoResponse_SignInfo
	uiSigns      []*widgets.Image
	uiButtonText *widgets.Paragraph
}

// Create a mock flipdot
func newMockFlipdotClient(signs []*flipdot.GetInfoResponse_SignInfo) mockFlipdotClient {
	// Create a widget for each image
	imageWidgets := []*widgets.Image{}
	previousHeight := 0
	for _, sign := range signs {
		// Create a new widget
		imgWidget := widgets.NewImage(nil)
		// Set the position of the image
		height := int(sign.Height) + 2 + previousHeight
		imgWidget.SetRect(0, previousHeight, int(sign.Width)+2, height)
		// Add the widget to the slice
		imageWidgets = append(imageWidgets, imgWidget)
		// Update previous height
		previousHeight = height
	}

	// Also create a paragraph to reveal button state
	buttonStateText := widgets.NewParagraph()

	return mockFlipdotClient{
		signConfig:   signs,
		uiSigns:      imageWidgets,
		uiButtonText: buttonStateText,
	}
}

// Mock the GetInfo response
func (m *mockFlipdotClient) GetInfo(ctx context.Context, in *flipdot.GetInfoRequest, opts ...grpc.CallOption) (*flipdot.GetInfoResponse, error) {
	// Return the baked-in mocked signs
	return &flipdot.GetInfoResponse{
		Signs: m.signs,
	}, nil
}

// Mock the Draw function
func (m *mockFlipdotClient) Draw(ctx context.Context, in *flipdot.DrawRequest, opts ...grpc.CallOption) (*flipdot.DrawResponse, error) {
	// Find the sign
	for i, sign := range m.signs {
		if sign.Name == in.Sign {
			// Draw the image
			img := unslice(*in.Image, sign.Width, sign.Height)
			// Draw the image to the terminal
			m.widgets[i].Image = img
			// Render the update
			termui.Render(m.widgets[i])
		}
	}
	return &flipdot.DrawResponse{}, nil
}

func (m *mockFlipdotClient) Test(ctx context.Context, in *flipdot.TestRequest, opts ...grpc.CallOption) (*flipdot.TestResponse, error) {
	// Return and do nothing
	return &flipdot.TestResponse{}, nil
}

func (m *mockFlipdotClient) Light(ctx context.Context, in *flipdot.LightRequest, opts ...grpc.CallOption) (*flipdot.LightResponse, error) {
	// Return and do nothing
	return &flipdot.LightResponse{}, nil
}

func unslice(imgIn flipdot.Image, width, height uint32) image.Image {
	// Create an image to hold the unpacked pixels
	imgOut := image.NewGray(image.Rect(0, 0, int(width), int(height)))
	// Iterate through the pixels
	for i, pixel := range imgIn.Data {
		row := i / int(width)
		col := i % int(width)
		var c color.Gray
		if pixel {
			c = color.Gray{255}
		} else {
			c = color.Gray{0}
		}
		imgOut.SetGray(col, row, c)
	}
	// Return the new image
	return imgOut
}

func uiHandlingLoop() {
	// Get the poll events channel
	uiEvents := termui.PollEvents()
	// Poll events until user quits
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "b": // press 'b' for button press

			}
		}
	}
}

// Create a mock flipdot for use in the application
func createMockFlipdotClient() flipdot.FlipdotClient {
	// Mock the signs
	mockSigns := []*flipdot.GetInfoResponse_SignInfo{
		&flipdot.GetInfoResponse_SignInfo{
			Name:   "top",
			Width:  84,
			Height: 7,
		},
		&flipdot.GetInfoResponse_SignInfo{
			Name:   "bottom",
			Width:  84,
			Height: 7,
		},
	}
	// Create the mock
	mock := newMockFlipdotClient(mockSigns)

	// Start a coroutine for checking for user input
	go func() {
		// Diable logging (we will be drawing)
		log.SetOutput(ioutil.Discard)
		// Initialise termui
		if err := termui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer termui.Close()
		// Listen for button presses, etc
		uiHandlingLoop()
	}()
	return &mock
}
