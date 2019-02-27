package cmd

import (
	"fmt"
	"testing"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

// Mock used to test the cli
var mock *flipdot.MockFlipdotClient
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

type drawImage struct {
	sign string
}

func RequestDrawImage(sign string) gomock.Matcher {
	return &drawImage{sign}
}
func (di *drawImage) Matches(x interface{}) bool {
	return x.(*flipdot.DrawRequest).Sign == di.sign
}
func (di *drawImage) String() string {
	return fmt.Sprintf("Sign should be %s", di.sign)
}

func createMock(t *testing.T) (*gomock.Controller, *flipdot.MockFlipdotClient) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock = flipdot.NewMockFlipdotClient(ctrl)
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

func TestText(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Construct signs
	top := flipdot.GetInfoResponse_SignInfo{
		Name: "top",
	}
	bottom := flipdot.GetInfoResponse_SignInfo{
		Name: "bottom",
	}
	info_response := flipdot.GetInfoResponse{
		Signs: []*flipdot.GetInfoResponse_SignInfo{&top, &bottom},
	}
	// Text command
	response := flipdot.DrawResponse{Error: &noError}
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("top")).Return(&response, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("bottom")).Return(&response, nil),
	)
	// Run the command
	rootCmd.SetArgs([]string{"text", "--font", "../Smirnof.ttf", "Hello my name is Sam. How's tricks?"})
	Execute(mockFactory)
}

func mockFactory(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
	return mock, nil, nil
}
