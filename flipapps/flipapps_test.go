package flipapps

import (
	context "context"
	reflect "reflect"
	"testing"
	"time"

	"github.com/briggySmalls/flipapp/flipdot"
	gomock "github.com/golang/mock/gomock"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

const (
	contextTimeoutS = 1
)

func TestGetInfo(t *testing.T) {
	ctrl, mock, flipapps := createTestObjects(t)
	defer ctrl.Finish()
	// Decide on some test values
	signs := []*flipdot.GetInfoResponse_SignInfo{
		&flipdot.GetInfoResponse_SignInfo{Name: "top", Width: 10, Height: 10},
		&flipdot.GetInfoResponse_SignInfo{Name: "bottom", Width: 20, Height: 20},
	}
	// Configure the mock
	mock.EXPECT().Signs().Return(signs)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	response, err := flipapps.GetInfo(ctx, &flipdot.GetInfoRequest{})
	// Assert the return values
	if err != nil {
		t.Fatal(err)
	}
	for _, sign := range response.Signs {
		if !reflect.DeepEqual(sign, sign) {
			t.Errorf("Signs don't match")
		}
	}
}

func TestMessageText(t *testing.T) {
	// Create mocks
	ctrl, mock, flipapps := createTestObjects(t)
	defer ctrl.Finish()
	// Create a channel to signal the test is complete
	complete := make(chan struct{})
	defer close(complete)
	// Configure the mock (calls 'done' when executed)
	mockAction := func(txt string, fnt font.Face, centre bool) {
		complete <- struct{}{}
	}
	mock.EXPECT().Text("test text", getTestFont(), false).Do(mockAction).Return(nil)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	_, err := flipapps.SendMessage(ctx, &MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	})
	// Assert the return values
	if err != nil {
		t.Fatal(err)
	}
	// Wait until the message is handled, or timeout
	select {
	case <-complete:
		// Completed successfully
		return
	case <-time.After(time.Second):
		// Timeout before we completed
		t.Fatal("Timeout before expected call")
	}
}

func getTestFont() font.Face {
	return inconsolata.Regular8x16
}

func createTestObjects(t *testing.T) (*gomock.Controller, *flipdot.MockFlipdot, FlipAppsServer) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock := flipdot.NewMockFlipdot(ctrl)
	// Create object under test
	flipapps := NewFlipappsServer(mock, getTestFont())
	return ctrl, mock, flipapps
}

// Get a simple context for sending requests via gRPC
func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeoutS*time.Second)
}
