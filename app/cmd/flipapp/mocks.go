package flipapp

import (
	"context"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/protos"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	rpio "github.com/stianeikeland/go-rpio/v4"
	"google.golang.org/grpc"
)

// Mock output pin
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

// Mock input pin
type mockInputPin struct {
	state  bool // Button state
	mux    sync.Mutex
	uiText *widgets.Paragraph
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
type mockUI struct {
	signConfig []*protos.GetInfoResponse_SignInfo
	uiSigns    []*widgets.Image
	buttonPin  mockInputPin
	ledPin     mockOutputPin
}

func newMockUI(signs []*protos.GetInfoResponse_SignInfo) mockUI {
	// Create an image widget for each sign
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

	// Create a text widget for the button
	buttonText := widgets.NewParagraph()
	buttonText.Text = "Button"
	buttonText.SetRect(0, previousHeight+1, 20, previousHeight+1+3)
	buttonText.TextStyle.Fg = termui.ColorWhite

	// Create a mockUI
	return mockUI{
		signConfig: signs,
		uiSigns:    imageWidgets,
		buttonPin:  mockInputPin{uiText: buttonText},
		ledPin:     mockOutputPin{uiText: buttonText},
	}
}

// Mock the GetInfo response
func (m *mockUI) GetInfo(ctx context.Context, in *protos.GetInfoRequest, opts ...grpc.CallOption) (*protos.GetInfoResponse, error) {
	// Return the baked-in mocked signs
	return &protos.GetInfoResponse{
		Signs: m.signConfig,
	}, nil
}

// Mock the Draw function
func (m *mockUI) Draw(ctx context.Context, in *protos.DrawRequest, opts ...grpc.CallOption) (*protos.DrawResponse, error) {
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
	return &protos.DrawResponse{}, nil
}

func (m *mockUI) Test(ctx context.Context, in *protos.TestRequest, opts ...grpc.CallOption) (*protos.TestResponse, error) {
	// Return and do nothing
	return &protos.TestResponse{}, nil
}

func (m *mockUI) Light(ctx context.Context, in *protos.LightRequest, opts ...grpc.CallOption) (*protos.LightResponse, error) {
	// Return and do nothing
	return &protos.LightResponse{}, nil
}

func unslice(imgIn protos.Image, width, height uint32) image.Image {
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

func (m *mockUI) ProcessEvents() {
	// Get the poll events channel
	uiEvents := termui.PollEvents()
	// Poll events until user quits
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "b": // Press button
				m.buttonPin.set(true)
				time.Sleep(time.Millisecond * 300)
				m.buttonPin.set(false)
				break
			}
		}
	}
}

// Create a mock flipdot for use in the application
func createMockUI() *mockUI {
	// Mock the signs
	signs := []*protos.GetInfoResponse_SignInfo{
		{
			Name:   "top",
			Width:  84,
			Height: 7,
		},
		{
			Name:   "bottom",
			Width:  84,
			Height: 7,
		},
	}
	// Create the UI
	ui := newMockUI(signs)

	// Run the UI loop
	go func() { // Start a coroutine for checking for user input
		// Diable logging (we will be drawing)
		log.SetOutput(ioutil.Discard)
		// Initialise termui
		if err := termui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		// Render all the objects
		for _, obj := range ui.uiSigns {
			termui.Render(obj)
		}
		termui.Render(ui.ledPin.uiText)
		defer termui.Close()
		// Listen for button presses, etc
		ui.ProcessEvents()
	}()

	// Return the ui
	return &ui
}
