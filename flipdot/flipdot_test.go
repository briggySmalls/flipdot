package flipdot

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

type testAction struct{ action TestRequest_Action }

func RequestTestAction(action TestRequest_Action) gomock.Matcher {
	return &testAction{action}
}
func (ta *testAction) Matches(x interface{}) bool {
	return x.(*TestRequest).Action == ta.action
}
func (ta *testAction) String() string {
	return fmt.Sprintf("Action is %s", ta.action.String())
}

type lightStatus struct{ status LightRequest_Status }

func RequestLightStatus(status LightRequest_Status) gomock.Matcher {
	return &lightStatus{status}
}
func (ls *lightStatus) Matches(x interface{}) bool {
	return x.(*LightRequest).Status == ls.status
}
func (ls *lightStatus) String() string {
	return fmt.Sprintf("Status is %s", ls.status.String())
}

type drawImage struct{ sign string }

func RequestDrawImage(sign string) gomock.Matcher {
	return &drawImage{sign}
}
func (di *drawImage) Matches(x interface{}) bool {
	return x.(*DrawRequest).Sign == di.sign
}
func (di *drawImage) String() string {
	return fmt.Sprintf("Sign should be %s", di.sign)
}

func TestCreate(t *testing.T) {
	// Create a mock
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure mock
	response := getStandardSignsResponse()
	mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&response, nil)
	// Create the flipdot instance
	_, err := NewFlipdot(mock)
	failOnError(err, t)
}

// Test sending the start test call
func TestTestStart(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	test_response := TestResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Test(gomock.Any(), RequestTestAction(TestRequest_START)).Return(&test_response, nil),
	)
	// Run the test
	runTest(func(f Flipdot) error {
		return f.TestStart()
	}, mock, t)
}

// Test sending the stop test call
func TestTestStop(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	test_response := TestResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Test(gomock.Any(), RequestTestAction(TestRequest_STOP)).Return(&test_response, nil),
	)
	// Run the test
	runTest(func(f Flipdot) error {
		return f.TestStop()
	}, mock, t)
}

// Test sending the lights on call
func TestLightsOn(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'On' command
	light_response := LightResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Light(gomock.Any(), RequestLightStatus(LightRequest_ON)).Return(&light_response, nil),
	)
	// Run the test
	runTest(func(f Flipdot) error {
		return f.LightOn()
	}, mock, t)
}

// Test sending the lights off call
func TestLightsOff(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// 'Off' command
	light_response := LightResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Light(gomock.Any(), RequestLightStatus(LightRequest_OFF)).Return(&light_response, nil),
	)
	// Run the test
	runTest(func(f Flipdot) error {
		return f.LightOff()
	}, mock, t)
}

// Test initialising with a client with incompatible signs
func TestDifferentSignsCaught(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Construct signs with different dimensions
	sign_a := GetInfoResponse_SignInfo{Name: "a", Width: 1, Height: 1}
	sign_b := GetInfoResponse_SignInfo{Name: "b", Width: 1, Height: 2}
	info_response := GetInfoResponse{Signs: []*GetInfoResponse_SignInfo{&sign_a, &sign_b}}
	// Configure the mock
	mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil)
	// Create a new flipdot
	_, err := NewFlipdot(mock)
	// Confirm there was an error
	if err == nil {
		t.Errorf("Incompatible signs not detected")
	}
}

// Test building text into images and sending them
func TestText(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	text_response := DrawResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("top")).Return(&text_response, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("bottom")).Return(&text_response, nil),
	)
	// Run the test
	runTest(func(f Flipdot) error {
		return f.Text("Hello my name is Sam. How's tricks?", getFont())
	}, mock, t)
}

// Helper function to create a mock FlipdotClient
func createMock(t *testing.T) (*gomock.Controller, *MockFlipdotClient) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock := NewMockFlipdotClient(ctrl)
	return ctrl, mock
}

// Helper function to create a Flipdot and run a test function
func runTest(fn func(f Flipdot) error, mock *MockFlipdotClient, t *testing.T) {
	// Create a flipdot
	f, err := NewFlipdot(mock)
	failOnError(err, t)
	// Run the command
	err = fn(f)
	failOnError(err, t)
}

// Helper function to get a font face
func getFont() (font font.Face) {
	return inconsolata.Regular8x16
}

func getStandardSignsResponse() (response GetInfoResponse) {
	// Construct signs
	top := GetInfoResponse_SignInfo{Name: "top", Width: 84, Height: 17}
	bottom := GetInfoResponse_SignInfo{Name: "bottom", Width: 84, Height: 17}
	response = GetInfoResponse{Signs: []*GetInfoResponse_SignInfo{&top, &bottom}}
	return
}

// Helper function to stop a test in the event of an error
func failOnError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
