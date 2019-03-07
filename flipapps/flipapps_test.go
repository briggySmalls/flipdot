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
		t.Errorf("GetInfo returned an error %s", err)
	}
	for _, sign := range response.Signs {
		if !reflect.DeepEqual(sign, sign) {
			t.Errorf("Signs don't match")
		}
	}
}

func TestMessageText(t *testing.T) {
	ctrl, mock, flipapps := createTestObjects(t)
	defer ctrl.Finish()
	// Configure the mock
	mock.EXPECT().Text("test text", getTestFont()).Return(nil)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	_, err := flipapps.SendMessage(ctx, &MessageRequest{
		From:    "briggySmalls",
		Payload: &MessageRequest_Text{"test text"},
	})
	// Assert the return values
	if err != nil {
		t.Errorf("GetInfo returned an error %s", err)
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
