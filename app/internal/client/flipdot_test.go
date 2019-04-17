package client

import (
	"fmt"
	reflect "reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

const (
	frameDuration = time.Millisecond
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

type drawImage struct {
	sign         string
	image        []bool
	isSignMatch  bool
	isImageMatch bool
}

func RequestDrawImage(sign string, image []bool) gomock.Matcher {
	return &drawImage{sign: sign, image: image}
}
func (di *drawImage) Matches(x interface{}) bool {
	di.isSignMatch = di.sign == x.(*DrawRequest).Sign
	di.isImageMatch = di.image == nil || reflect.DeepEqual(x.(*DrawRequest).Image.Data, di.image)
	return di.isSignMatch && di.isImageMatch
}
func (di *drawImage) String() string {
	if !di.isImageMatch {
		return fmt.Sprintf("Image should match")
	} else if !di.isSignMatch {
		return fmt.Sprintf("Sign should be '%s'", di.sign)
	} else {
		panic("String() called but no detected failure")
	}
}

func TestCreate(t *testing.T) {
	// Create a mock
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure mock
	response := getStandardSignsResponse()
	mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&response, nil)
	// Create the flipdot instance
	_, err := NewFlipdot(mock, frameDuration)
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
	_, err := NewFlipdot(mock, frameDuration)
	// Confirm there was an error
	if err == nil {
		t.Errorf("Incompatible signs not detected")
	}
}

func TestDraw(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	drawResponse := DrawResponse{}
	infoResponse := getStandardSignsResponse()
	// Make some mock images
	falseImageData := make([]bool, infoResponse.Signs[0].Width*infoResponse.Signs[0].Height)
	topImageData := make([]bool, infoResponse.Signs[0].Width*infoResponse.Signs[0].Height)
	topImageData[0] = true
	bottomImageData := make([]bool, infoResponse.Signs[0].Width*infoResponse.Signs[0].Height)
	bottomImageData[1] = true
	// Expect the mock images
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&infoResponse, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("top", topImageData)).Return(&drawResponse, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("bottom", bottomImageData)).Return(&drawResponse, nil),
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("top", topImageData)).Return(&drawResponse, nil),
		// Final image is false because it is an 'empty' end-of-frame
		mock.EXPECT().Draw(gomock.Any(), RequestDrawImage("bottom", falseImageData)).Return(&drawResponse, nil),
	)
	// Create a couple of images
	images := []*Image{
		&Image{Data: topImageData},
		&Image{Data: bottomImageData},
		&Image{Data: topImageData},
	}
	// Run the test
	runTest(func(f Flipdot) error {
		return f.Draw(images, false)
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
	f, err := NewFlipdot(mock, frameDuration)
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
