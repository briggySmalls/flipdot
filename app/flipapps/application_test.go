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
	ctrl, mock, app := createAppTestObjects(t, time.Millisecond)
	defer ctrl.Finish()
	// Start the app running
	go app.Run()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	// Configure the mock (calls 'done' when executed)
	mockAction := func(txt string, fnt font.Face, centre bool) {
		// Assert that the string is as expected
		_, err := time.Parse("Mon 2 Jan\n03:04 PM", txt)
		if err != nil {
			t.Fatal(err)
		}
		// Finish up
		textWritten <- struct{}{}
	}
	mock.EXPECT().Text(gomock.Any(), getTestFont(), true).Do(mockAction).Return(nil)
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully, stop the app
		close(app.MessageQueue)
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func TestMessageText(t *testing.T) {
	// Create mocks
	ctrl, mock, app := createAppTestObjects(t, time.Hour)
	defer ctrl.Finish()
	// Start the app running
	go app.Run()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	// Configure the mock (calls 'done' when executed)
	mockAction := func(txt string, fnt font.Face, centre bool) {
		textWritten <- struct{}{}
	}
	mock.EXPECT().Text("test text", getTestFont(), false).Do(mockAction).Return(nil)
	// Send the message
	message := MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	app.MessageQueue <- message
	// Wait until the message is handled, or timeout
	select {
	case <-textWritten:
		// Completed successfully, stop the app
		close(app.MessageQueue)
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func createAppTestObjects(t *testing.T, tickTime time.Duration) (*gomock.Controller, *flipdot.MockFlipdot, application) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock := flipdot.NewMockFlipdot(ctrl)
	// Create object under test
	app := NewApplication(mock, tickTime, getTestFont())
	return ctrl, mock, app
}

func getTestFont() font.Face {
	return inconsolata.Regular8x16
}
