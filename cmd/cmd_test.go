package cmd

import (
	"fmt"
	"testing"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

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
	mock := flipdot.NewMockFlipdotClient(ctrl)
	return ctrl, mock
}

func TestTestStart(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	nilErr := getNoError()
	response := flipdot.TestResponse{Error: &nilErr}
	mock.EXPECT().Test(gomock.Any(), RequestTestAction(flipdot.TestRequest_START)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"test", "start"})
	Execute(func(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestTestStop(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	nilErr := getNoError()
	response := flipdot.TestResponse{Error: &nilErr}
	mock.EXPECT().Test(gomock.Any(), RequestTestAction(flipdot.TestRequest_STOP)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"test", "stop"})
	Execute(func(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestLightsOn(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'On' command
	nilErr := getNoError()
	response := flipdot.LightResponse{Error: &nilErr}
	mock.EXPECT().Light(gomock.Any(), RequestLightStatus(flipdot.LightRequest_ON)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"light", "on"})
	Execute(func(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestLightsOff(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'Off' command
	nilErr := getNoError()
	response := flipdot.LightResponse{Error: &nilErr}
	mock.EXPECT().Light(gomock.Any(), RequestLightStatus(flipdot.LightRequest_OFF)).Return(&response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"light", "off"})
	Execute(func(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestDifferentSignsCaught(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Construct signs with different dimensions
	sign_a := flipdot.GetInfoResponse_SignInfo{
		Name:   "a",
		Width:  1,
		Height: 1,
	}
	sign_b := flipdot.GetInfoResponse_SignInfo{
		Name:   "b",
		Width:  1,
		Height: 2,
	}
	info_response := flipdot.GetInfoResponse{Signs: []*flipdot.GetInfoResponse_SignInfo{&sign_a, &sign_b}}
	// Prepare for failure
	defer func() {
		if err := recover(); err != nil {
			// We paniced
		}
	}()
	mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil)
	// Run the command
	rootCmd.SetArgs([]string{"text", "--font", "../Smirnof.ttf", "Hello my name is Sam. How's tricks?"})
	Execute(func(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
	// We didn't panic
	t.Fail()
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
	nilErr := getNoError()
	response := flipdot.DrawResponse{Error: &nilErr}
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("top")).Return(&response, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("bottom")).Return(&response, nil),
	)
	// Run the command
	rootCmd.SetArgs([]string{"text", "--font", "../Smirnof.ttf", "Hello my name is Sam. How's tricks?"})
	Execute(func(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func getNoError() flipdot.Error {
	return flipdot.Error{Code: 0}
}
