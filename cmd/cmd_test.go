package cmd

import (
	"testing"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

func createMock(t *testing.T) (*gomock.Controller, *flipdot.MockFlipdot) {
	// Create a mock
	ctrl := gomock.NewController(t)
	mock := flipdot.NewMockFlipdot(ctrl)
	return ctrl, mock
}

func TestTestStart(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	mock.EXPECT().TestStart()
	// Run the command
	rootCmd.SetArgs([]string{"test", "start"})
	Execute(func(port uint) (flipdot.Flipdot, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestTestStop(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	mock.EXPECT().TestStop()
	// Run the command
	rootCmd.SetArgs([]string{"test", "stop"})
	Execute(func(port uint) (flipdot.Flipdot, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestLightsOn(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	mock.EXPECT().LightOn()
	// Run the command
	rootCmd.SetArgs([]string{"light", "on"})
	Execute(func(port uint) (flipdot.Flipdot, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}

func TestLightsOff(t *testing.T) {
	ctrl, mock := createMock(t)
	defer ctrl.Finish()
	// Configure the mock
	mock.EXPECT().LightOff()
	// Run the command
	rootCmd.SetArgs([]string{"light", "off"})
	Execute(func(port uint) (flipdot.Flipdot, *grpc.ClientConn, error) {
		return mock, nil, nil
	})
}
