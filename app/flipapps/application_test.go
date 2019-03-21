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
	ctrl, fakeFlipdot, fakeBm, app := createAppTestObjects(t, time.Millisecond)
	defer ctrl.Finish()
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
	fakeBm.EXPECT().GetChannel()
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

func TestMessageTextQueued(t *testing.T) {
	// Create mocks
	ctrl, _, fakeBm, app := createAppTestObjects(t, time.Hour)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	messageAdded := make(chan struct{})
	defer close(messageAdded)
	// Configure mock to expect a call to activate button
	fakeBm.EXPECT().GetChannel()
	fakeBm.EXPECT().SetState(Active).Do(func(state State) {
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

func TestMessageTextSent(t *testing.T) {
	// Create mocks
	ctrl, fakeFlipdot, fakeBm, app := createAppTestObjects(t, time.Hour)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	textWritten := make(chan struct{})
	defer close(textWritten)
	defer close(app.MessagesIn)
	// Create a channel to signal a button press
	buttonPress := make(chan struct{})
	// Configure the mocks
	fakeBm.EXPECT().GetChannel().Return(buttonPress)
	fakeBm.EXPECT().SetState(Inactive)
	fakeBm.EXPECT().SetState(Active)
	fakeFlipdot.EXPECT().Text("test text", getTestFont(), false).Do(func(txt string, fnt font.Face, centre bool) {
		textWritten <- struct{}{}
	}).Return(nil)

	// Send a message to start the test
	app.MessagesIn <- MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	}
	// Then send a button press
	buttonPress <- struct{}{}
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
