package client

import (
	"fmt"
	reflect "reflect"
	"testing"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/protos"
	"github.com/golang/mock/gomock"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
)

const (
	frameDuration = time.Millisecond
)

type testAction struct{ action protos.TestRequest_Action }

func RequestTestAction(action protos.TestRequest_Action) gomock.Matcher {
	return &testAction{action}
}
func (ta *testAction) Matches(x interface{}) bool {
	return x.(*protos.TestRequest).Action == ta.action
}
func (ta *testAction) String() string {
	return fmt.Sprintf("Action is %s", ta.action.String())
}

type lightStatus struct{ status protos.LightRequest_Status }

func RequestLightStatus(status protos.LightRequest_Status) gomock.Matcher {
	return &lightStatus{status}
}
func (ls *lightStatus) Matches(x interface{}) bool {
	return x.(*protos.LightRequest).Status == ls.status
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
	di.isSignMatch = di.sign == x.(*protos.DrawRequest).Sign
	di.isImageMatch = di.image == nil || reflect.DeepEqual(x.(*protos.DrawRequest).Image.Data, di.image)
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
	test_response := protos.TestResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Test(gomock.Any(), RequestTestAction(protos.TestRequest_START)).Return(&test_response, nil),
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
	test_response := protos.TestResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Test(gomock.Any(), RequestTestAction(protos.TestRequest_STOP)).Return(&test_response, nil),
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
	light_response := protos.LightResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Light(gomock.Any(), RequestLightStatus(protos.LightRequest_ON)).Return(&light_response, nil),
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
	light_response := protos.LightResponse{}
	info_response := getStandardSignsResponse()
	gomock.InOrder(
		mock.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&info_response, nil),
		mock.EXPECT().Light(gomock.Any(), RequestLightStatus(protos.LightRequest_OFF)).Return(&light_response, nil),
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
	sign_a := protos.GetInfoResponse_SignInfo{Name: "a", Width: 1, Height: 1}
	sign_b := protos.GetInfoResponse_SignInfo{Name: "b", Width: 1, Height: 2}
	info_response := protos.GetInfoResponse{Signs: []*protos.GetInfoResponse_SignInfo{&sign_a, &sign_b}}
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
	drawResponse := protos.DrawResponse{}
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
	images := []*protos.Image{
		{Data: topImageData},
		{Data: bottomImageData},
		{Data: topImageData},
	}
	// Run the test
	runTest(func(f Flipdot) error {
		return f.Draw(images, false)
	}, mock, t)
}

// Helper function to create a mock FlipdotClient
func createMock(t *testing.T) (*gomock.Controller, *protos.MockDriverClient) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock := protos.NewMockDriverClient(ctrl)
	return ctrl, mock
}

// Helper function to create a Flipdot and run a test function
func runTest(fn func(f Flipdot) error, mock *protos.MockDriverClient, t *testing.T) {
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

func getStandardSignsResponse() (response protos.GetInfoResponse) {
	// Construct signs
	top := protos.GetInfoResponse_SignInfo{Name: "top", Width: 84, Height: 17}
	bottom := protos.GetInfoResponse_SignInfo{Name: "bottom", Width: 84, Height: 17}
	response = protos.GetInfoResponse{Signs: []*protos.GetInfoResponse_SignInfo{&top, &bottom}}
	return
}

// Helper function to stop a test in the event of an error
func failOnError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
