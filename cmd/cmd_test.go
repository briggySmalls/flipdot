package cmd

import (
	"testing"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/briggySmalls/flipcli/mock_flipdot"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

// Mock used to test the cli
var mock *mock_flipdot.MockFlipdotClient

func TestStartTestSigns(t *testing.T) {
	// Create a mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock = mock_flipdot.NewMockFlipdotClient(ctrl)
	// Configure the mock
	no_error := flipdot.Error{Code: 0}
	response := flipdot.StartTestResponse{Error: &no_error}
	mock.EXPECT().StartTest(gomock.Any(), gomock.Any()).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"--port", "8000", "test", "start"})
	Execute(mockFactory)
}

func mockFactory(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
	return mock, nil, nil
}
