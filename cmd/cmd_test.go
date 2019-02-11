package cmd

import (
	"fmt"
	"testing"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/briggySmalls/flipcli/mock_flipdot"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

// Mock used to test the cli
var mock *mock_flipdot.MockFlipdotClient
var noError = flipdot.Error{Code: 0}

type testAction struct{ action flipdot.TestRequest_Action }

func RequestTestAction(action flipdot.TestRequest_Action) gomock.Matcher {
	return &testAction{action}
}
func (ta *testAction) Matches(x interface{}) bool {
	return x.(*flipdot.TestRequest).Action == ta.action
}
func (ta *testAction) String() string {
	return fmt.Sprintf("Action is %s", ta.action.String())
}

type lightStatus struct{ status flipdot.LightRequest_Status }

func RequestLightStatus(status flipdot.LightRequest_Status) gomock.Matcher {
	return &lightStatus{status}
}
func (ls *lightStatus) Matches(x interface{}) bool {
	return x.(*flipdot.LightRequest).Status == ls.status
}
func (ls *lightStatus) String() string {
	return fmt.Sprintf("Status is %s", ls.status.String())
}

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
	response := flipdot.TestResponse{Error: &noError}
	mock.EXPECT().Test(gomock.Any(), RequestTestAction(flipdot.TestRequest_START)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"test", "start"})
	Execute(mockFactory)
}

func TestTestStop(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	response := flipdot.TestResponse{Error: &noError}
	mock.EXPECT().Test(gomock.Any(), RequestTestAction(flipdot.TestRequest_STOP)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"test", "stop"})
	Execute(mockFactory)
}

func TestLightsOn(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'On' command
	response := flipdot.LightResponse{Error: &noError}
	mock.EXPECT().Light(gomock.Any(), RequestLightStatus(flipdot.LightRequest_ON)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"light", "on"})
	Execute(mockFactory)
}

func TestLightsOff(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'Off' command
	response := flipdot.LightResponse{Error: &noError}
	mock.EXPECT().Light(gomock.Any(), RequestLightStatus(flipdot.LightRequest_OFF)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"light", "off"})
	Execute(mockFactory)
}

func mockFactory(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
	return mock, nil, nil
}
