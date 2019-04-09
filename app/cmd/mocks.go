package cmd

import (
	"context"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"sync"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/stianeikeland/go-rpio"
	"google.golang.org/grpc"
)

// Mock output pin
type mockInputPin struct {
	state  bool // Button state
	mux    sync.Mutex
	uiText *widgets.Paragraph
}

type mockOutputPin struct {
	uiText *widgets.Paragraph
	state  bool // Button state
}

func (p *mockOutputPin) High()   { p.update(true) }
func (p *mockOutputPin) Low()    { p.update(false) }
func (p *mockOutputPin) Toggle() { p.update(!p.state) }
func (p *mockOutputPin) update(state bool) {
	// Update state
	p.state = state
	// Update paragraph colour
	color := &p.uiText.TextStyle.Fg
	if state {
		*color = termui.ColorRed
	} else {
		*color = termui.ColorWhite
	}
	// Render the update
	termui.Render(p.uiText)
}

func (p *mockInputPin) Read() rpio.State {
	p.mux.Lock()
	defer p.mux.Unlock()
	if p.state {
		return rpio.High
	}
	return rpio.Low
}

func (p *mockInputPin) set(state bool) {
	p.mux.Lock()
	p.state = state
	p.uiText.Border = !state
	// Render the update
	termui.Render(p.uiText)
	p.mux.Unlock()
}

// Mock flipdot
type mockFlipdotClient struct {
	signConfig []*flipdot.GetInfoResponse_SignInfo
	uiSigns    []*widgets.Image
	buttonPin  mockInputPin
	ledPin     mockOutputPin
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

	// Create some text for the button
	buttonText := widgets.NewParagraph()
	buttonText.Text = "Button"
	buttonText.SetRect(0, previousHeight+1, 20, previousHeight+1+3)
	buttonText.TextStyle.Fg = termui.ColorWhite

	return mockFlipdotClient{
		signConfig: signs,
		uiSigns:    imageWidgets,
		buttonPin:  mockInputPin{uiText: buttonText},
		ledPin:     mockOutputPin{uiText: buttonText},
	}
}

// Mock the GetInfo response
func (m *mockFlipdotClient) GetInfo(ctx context.Context, in *flipdot.GetInfoRequest, opts ...grpc.CallOption) (*flipdot.GetInfoResponse, error) {
	// Return the baked-in mocked signs
	return &flipdot.GetInfoResponse{
		Signs: m.signConfig,
	}, nil
}

// Mock the Draw function
func (m *mockFlipdotClient) Draw(ctx context.Context, in *flipdot.DrawRequest, opts ...grpc.CallOption) (*flipdot.DrawResponse, error) {
	// Find the sign
	for i, sign := range m.signConfig {
		if sign.Name == in.Sign {
			// Draw the image
			img := unslice(*in.Image, sign.Width, sign.Height)
			// Draw the image to the terminal
			m.uiSigns[i].Image = img
			// Render the update
			termui.Render(m.uiSigns[i])
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

func (m *mockFlipdotClient) uiHandlingLoop() {
	// Get the poll events channel
	uiEvents := termui.PollEvents()
	// Poll events until user quits
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "<Down>": // Press button
				m.buttonPin.set(true)
				break
			case "<Up>": // Release button
				m.buttonPin.set(false)
			}
		}
	}
}

// Create a mock flipdot for use in the application
func createMockFlipdotClient() *mockFlipdotClient {
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
		mock.uiHandlingLoop()
	}()
	return &mock
}