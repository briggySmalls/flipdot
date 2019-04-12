package flipapps

import (
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
	defer close(app.GetMessagesChannel())
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	// Configure imager mock to expect request to build images
	// Return two 'images' to be displayed
	fakeImager.EXPECT().Clock(gomock.Any(), false).Return([]*flipdot.Image{
		{Data: make([]bool, 10)},
		{Data: make([]bool, 10)},
	}, nil)
	// Configure the mock (calls 'done' when executed)
	mockAction := func(images []*flipdot.Image, isWait bool) {
		// Assert that the images are as expected
		if len(images) != 2 {
			t.Errorf("Unexpected number of images: %d", len(images))
		}
		// Finish up
		textWritten <- struct{}{}
	}
	fakeBm.EXPECT().GetChannel()
	fakeFlipdot.EXPECT().Draw(gomock.Any(), false).Do(mockAction).Return(nil)
	// Run
	go app.Run(time.Millisecond)
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestMessageTextQueued(t *testing.T) {
	// Create mocks
	ctrl, fakeFlipdot, fakeBm, fakeImager, app := createAppTestObjects(t)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	messageAdded := make(chan struct{})
	defer close(messageAdded)
	// Configure mock to expect a call to activate button
	gomock.InOrder(
		fakeBm.EXPECT().GetChannel(),                   // Expect setup to fetch channel (before loop)
		fakeImager.EXPECT().Clock(gomock.Any(), false), // Expect clock image to be built (before loop)
		fakeFlipdot.EXPECT().Draw(gomock.Any(), false), // Expect clock image to be drawn (before loop)
		fakeBm.EXPECT().SetState(Active),               // Expect button to be activated
		fakeImager.EXPECT().Clock(gomock.Any(), true),  // Expect clock image to be built
		fakeFlipdot.EXPECT().Draw(gomock.Any(), false).Do(func(interface{}, bool) {
			// We are done testing
			messageAdded <- struct{}{}
		}), // Expect clock images to be sent
	)
	// Run
	go app.Run(time.Millisecond)
	// Send the message
	messagesIn := app.GetMessagesChannel()
	defer close(messagesIn)
	messagesIn <- MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	// Wait until the message is handled, or timeout
	select {
	case <-messageAdded:
		// Completed successfully, stop the app
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
	activated := make(chan struct{})
	defer close(activated)
	// Get channel to pass messages through
	messagesIn := app.GetMessagesChannel()
	defer close(messagesIn)
	// Configure startup mocks
	buttonPress := make(chan struct{}) // Create a channel to signal a button press
	gomock.InOrder(
		fakeBm.EXPECT().GetChannel().Return(buttonPress), // Pass button press channel to app, when asked
		fakeImager.EXPECT().Clock(gomock.Any(), false),   // Expect clock image to be built
		fakeFlipdot.EXPECT().Draw(gomock.Any(), false),   // Expect clock images to be sent
		fakeBm.EXPECT().SetState(Active).Do(func(interface{}) {
			// Signal to main thread that button was activated
			// Note: Can't message buttonPress in this callback as we get deadlock
			activated <- struct{}{}
		}), // Expect button to be activated after receiving message,
		fakeImager.EXPECT().Clock(gomock.Any(), true),  // Expect clock image to be built
		fakeFlipdot.EXPECT().Draw(gomock.Any(), false), // Expect clock images to be sent
		fakeBm.EXPECT().SetState(Inactive),             // Expect dectivate before drawing message
		fakeImager.EXPECT().Message("briggySmalls", "test text").Return([]*flipdot.Image{ // Expect constructing message images
			{Data: make([]bool, 10)},
			{Data: make([]bool, 10)},
			{Data: make([]bool, 10)},
			{Data: make([]bool, 10)},
		}, nil),
		fakeFlipdot.EXPECT().Draw(gomock.Any(), true).Do(func(images []*flipdot.Image, isWait bool) { // Expect draw message images
			// Check image
			if len(images) != 4 {
				t.Errorf("Unexpected number of images: %d", len(images))
			}
			// signal we are done
			textWritten <- struct{}{}
		}).Return(nil),
	)
	// Start the app
	go app.Run(time.Hour)
	// Send a message to start the test (note: we don't assert as we check this in previous test)
	messagesIn <- MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	// Wait until the message is handled, or timeout
	for {
		select {
		case <-activated:
			// Button is now active, so press button
			buttonPress <- struct{}{}
		case <-textWritten:
			// Completed successfully, stop the app
			return
		case <-time.After(time.Second * 5):
			// Timeout before we completed
			t.Fatal("Timeout before expected call")
		}
	}
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
