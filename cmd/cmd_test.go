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
var no_error = flipdot.Error{Code: 0}

func createMock(t *testing.T) (*gomock.Controller, *mock_flipdot.MockFlipdotClient) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock = mock_flipdot.NewMockFlipdotClient(ctrl)
	return ctrl, mock
}

func TestTestStart(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	response := flipdot.TestResponse{Error: &no_error}
	request := flipdot.TestRequest{Action: flipdot.TestRequest_START}
	mock.EXPECT().Test(gomock.Any(), gomock.Eq(&request)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"test", "start"})
	Execute(mockFactory)
}

func TestTestStop(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	response := flipdot.TestResponse{Error: &no_error}
	request := flipdot.TestRequest{Action: flipdot.TestRequest_STOP}
	mock.EXPECT().Test(gomock.Any(), gomock.Eq(request)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"test", "stop"})
	Execute(mockFactory)
}

func TestLightsOn(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'On' command
	response := flipdot.LightResponse{Error: &no_error}
	request := flipdot.LightRequest{Status: flipdot.LightRequest_ON}
	mock.EXPECT().Light(gomock.Any(), gomock.Eq(request)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"light", "on"})
	Execute(mockFactory)
}

func TestLightsOff(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'Off' command
	request := flipdot.LightRequest{Status: flipdot.LightRequest_OFF}
	response := flipdot.LightResponse{Error: &no_error}
	mock.EXPECT().Light(gomock.Any(), gomock.Eq(request)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"light", "off"})
	Execute(mockFactory)
}

func mockFactory(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
	return mock, nil, nil
}
