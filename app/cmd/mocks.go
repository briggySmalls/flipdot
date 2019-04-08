package cmd

import (
	"image"
	"image/color"
	"log"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"github.com/stianeikeland/go-rpio"
)

// Mock output pin
type mockOutputPin struct{}

func (p *mockOutputPin) High()   {}
func (p *mockOutputPin) Low()    {}
func (p *mockOutputPin) Toggle() {}

// Mock trigger pin
type mockTriggerPin struct{}

func (p *mockTriggerPin) EdgeDetected() bool { return false }
func (p *mockTriggerPin) Read() rpio.State   { return rpio.Low }

// Mock flipdot
type mockFlipdot struct {
	signs   []*flipdot.GetInfoResponse_SignInfo
	widgets []*widgets.Image
}

// NewMockFlipdot Create a mock flipdot
func NewMockFlipdot(signs []*flipdot.GetInfoResponse_SignInfo) flipdot.Flipdot {
	// Create a widget for each image
	imageWidgets := []*widgets.Image{}
	previousHeight := 0
	for _, sign := range signs {
		// Create a new widget
		imgWidget := widgets.NewImage(nil)
		// Set the position of the image
		height := int(sign.Height) + previousHeight
		imgWidget.SetRect(0, previousHeight, int(sign.Width), height)
		// Add the widget to the slice
		imageWidgets = append(imageWidgets, imgWidget)
		// Update previous height
		previousHeight = height
	}

	return &mockFlipdot{
		signs:   signs,
		widgets: imageWidgets,
	}
}

// Mock the GetInfo response
func (m *mockFlipdot) Signs() []*flipdot.GetInfoResponse_SignInfo {
	return m.signs
}

// Mock the Size response
func (m *mockFlipdot) Size() (width, height uint) {
	return uint(m.signs[0].Width), uint(m.signs[0].Height)
}

func (m *mockFlipdot) LightOn() error   { return nil }
func (m *mockFlipdot) LightOff() error  { return nil }
func (m *mockFlipdot) TestStart() error { return nil }
func (m *mockFlipdot) TestStop() error  { return nil }

// Mock drawing an image
func (m *mockFlipdot) Draw(images []*flipdot.Image) error {
	// Draw the images
	for i, image := range images {
		log.Println("Drawing image...")
		// Unslice the image
		img := m.unslice(*image)
		// Draw the image to the terminal
		m.widgets[i].Image = img
		termui.Render(m.widgets[i])
	}
	// Never error
	return nil
}

func (m *mockFlipdot) unslice(imgIn flipdot.Image) image.Image {
	// Create an image to hold the unpacked pixels
	width, height := m.Size()
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

// Create a mock flipdot for use in the application
func createMockFlipdot() flipdot.Flipdot {
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
	mock := NewMockFlipdot(mockSigns)

	return mock
}
