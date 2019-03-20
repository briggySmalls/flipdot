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
	ctrl, fakeFlipdot, _, app := createAppTestObjects(t, time.Millisecond)
	defer ctrl.Finish()
	// Start the app running
	go app.Run()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	// Configure the mock (calls 'done' when executed)
	mockAction := func(txt string, fnt font.Face, centre bool) {
		// Assert that the string is as expected
		_, err := time.Parse("Mon 2 Jan\n3:04 PM", txt)
		if err != nil {
			t.Fatal(err)
		}
		// Finish up
		textWritten <- struct{}{}
	}
	fakeFlipdot.EXPECT().Text(gomock.Any(), getTestFont(), true).Do(mockAction).Return(nil)
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully, stop the app
		close(app.MessagesIn)
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestLightOn(t *testing.T) {
	// Create mocks
	ctrl, _, fakeBm, app := createAppTestObjects(t, time.Hour)
	defer ctrl.Finish()
	// Start the app running
	go app.Run()
	// Create a channel to signal the test is complete
	messageAdded := make(chan struct{})
	defer close(messageAdded)
	// Configure mock to expect a call to light button LED
	fakeBm.EXPECT().WriteLed(true).Do(func(status bool) {
		messageAdded <- struct{}{}
	})
	// Send the message
	message := MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	app.MessagesIn <- message
	// Wait until the message is handled, or timeout
	select {
	case <-messageAdded:
		// Completed successfully, stop the app
		close(app.MessagesIn)
		return
	case <-time.After(time.Second * 5):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestMessageText(t *testing.T) {
	// Create mocks
	ctrl, fakeFlipdot, fakeBm, app := createAppTestObjects(t, time.Hour)
	defer ctrl.Finish()
	// Start the app running
	go app.Run()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	defer close(app.MessagesIn)
	defer close(app.ShowMessage)
	// Configure the mock (calls 'done' when executed)
	fakeBm.EXPECT().WriteLed(true)
	fakeFlipdot.EXPECT().Text("test text", getTestFont(), false).Do(func(txt string, fnt font.Face, centre bool) {
		textWritten <- struct{}{}
	}).Return(nil)
	// Send the message
	message := MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	app.MessagesIn <- message
	app.ShowMessage <- struct{}{}
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully, stop the app
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func createAppTestObjects(t *testing.T, tickTime time.Duration) (*gomock.Controller, *flipdot.MockFlipdot, *MockButtonManager, application) {
	// Create a mock
	ctrl := gomock.NewController(t)
	fakeFlipdot := flipdot.NewMockFlipdot(ctrl)
	fakeBm := NewMockButtonManager(ctrl)
	// Create object under test
	app := NewApplication(fakeFlipdot, fakeBm, tickTime, getTestFont())
	return ctrl, fakeFlipdot, fakeBm, app
}

func getTestFont() font.Face {
	return inconsolata.Regular8x16
}
