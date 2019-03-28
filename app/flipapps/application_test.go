package flipapps

import (
	"image"
	"image/color"
	"image/draw"
	reflect "reflect"
	"testing"
	"time"

	"github.com/briggySmalls/flipdot/app/flipdot"
	gomock "github.com/golang/mock/gomock"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

func TestTickText(t *testing.T) {
	// Create mocks
	ctrl, fakeFlipdot, fakeBm, fakeImager, app := createAppTestObjects(t)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	// Configure imager mock to expect request to build images
	// Return two 'images' to be displayed
	fakeImager.EXPECT().Clock(gomock.Any(), false).Return([]*flipdot.Image{
		&flipdot.Image{Data: make([]bool, 10)},
		&flipdot.Image{Data: make([]bool, 10)},
	}, nil)
	// Configure the mock (calls 'done' when executed)
	mockAction := func(images []*flipdot.Image) {
		// Assert that the images are as expected
		if len(images) != 2 {
			t.Errorf("Unexpected number of images: %d", len(images))
		}
		// Finish up
		textWritten <- struct{}{}
	}
	fakeBm.EXPECT().GetChannel()
	fakeFlipdot.EXPECT().Draw(gomock.Any()).Do(mockAction).Return(nil)
	// Run
	go app.Run(time.Millisecond)
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully, stop the app
		close(app.GetMessagesChannel())
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestMessageTextQueued(t *testing.T) {
	// Create mocks
	ctrl, _, fakeBm, _, app := createAppTestObjects(t)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	messageAdded := make(chan struct{})
	defer close(messageAdded)
	// Configure mock to expect a call to activate button
	fakeBm.EXPECT().GetChannel()
	fakeBm.EXPECT().SetState(Active).Do(func(state State) {
		messageAdded <- struct{}{}
	})
	// Run
	go app.Run(time.Millisecond)
	// Send the message
	messagesIn := app.GetMessagesChannel()
	messagesIn <- MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	// Wait until the message is handled, or timeout
	select {
	case <-messageAdded:
		// Completed successfully, stop the app
		close(messagesIn)
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestMessageTextSent(t *testing.T) {
	// Create mocks
	ctrl, fakeFlipdot, fakeBm, fakeImager, app := createAppTestObjects(t)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	messagesIn := app.GetMessagesChannel()
	defer close(messagesIn)
	// Configure startup mocks
	buttonPress := make(chan struct{}) // Create a channel to signal a button press
	fakeBm.EXPECT().GetChannel().Return(buttonPress)
	// Start the app
	app.Run(time.Hour)
	// Send a message to start the test (note: we don't assert as we check this in previous test)
	messagesIn <- MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	fakeBm.EXPECT().SetState(Inactive)                                         // Expect dectivate before drawing message
	fakeImager.EXPECT().Message("briggySmalls", "test text")                   // Expect first message
	fakeFlipdot.EXPECT().Draw(gomock.Any()).Do(func(images []*flipdot.Image) { // Expect draw first message
		// Check image
		if len(images) != 1 {
			t.Errorf("Unexpected number of images: %d", len(images))
		}
		// signal we are done
	}).Return(nil)
	// Then send a button press
	buttonPress <- struct{}{}
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully, stop the app
		return
	case <-time.After(time.Second * 5):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestSlice(t *testing.T) {
	// Prepare test table
	tables := []struct {
		input  image.Image
		output []bool
	}{
		{createTestImage(color.Gray{1}, image.Rect(0, 0, 2, 3)), []bool{true, true, true, true, true, true}},
		{createTestImage(color.Gray{0}, image.Rect(0, 0, 3, 3)), []bool{false, false, false, false, false, false, false, false, false}},
	}
	for _, table := range tables {
		slice := Slice(table.input)
		if !reflect.DeepEqual(slice, table.output) {
			t.Errorf("Image %s not sliced correctly", table.input)
		}
	}
}

func createTestImage(c color.Color, r image.Rectangle) image.Image {
	// Draw the image
	dst := image.NewGray(r)
	draw.Draw(dst, dst.Bounds(), image.NewUniform(c), image.Point{X: 0, Y: 0}, draw.Src)
	return dst
}

func createAppTestObjects(t *testing.T) (*gomock.Controller, *flipdot.MockFlipdot, *MockButtonManager, *MockImager, Application) {
	// Create a mock
	ctrl := gomock.NewController(t)
	fakeFlipdot := flipdot.NewMockFlipdot(ctrl)
	fakeBm := NewMockButtonManager(ctrl)
	fakeImager := NewMockImager(ctrl)
	// Create object under test
	app := NewApplication(fakeFlipdot, fakeBm, fakeImager)
	return ctrl, fakeFlipdot, fakeBm, fakeImager, app
}

func getTestFont() font.Face {
	return inconsolata.Regular8x16
}
