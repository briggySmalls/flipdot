package flipapps

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
)

func TestLightLed(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	// Create a button manager
	bm := NewButtonManager(fakeButtonPin, fakeLedPin)
	// Configure mock to expect calls
	gomock.InOrder(
		fakeLedPin.EXPECT().High(),
		fakeLedPin.EXPECT().Low(),
	)
	// Simulate a call to turn on the LED
	bm.WriteLed(true)
	bm.WriteLed(false)
}

func TestButtonPressed(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	// Configure mock to expect calls
	gomock.InOrder(
		fakeButtonPin.EXPECT().EdgeDetected().Times(100).Return(false),
		fakeButtonPin.EXPECT().EdgeDetected().Times(1).Return(true),
		fakeButtonPin.EXPECT().EdgeDetected().AnyTimes().Return(false),
	)
	// Create a button manager and start the app
	bm := NewButtonManager(fakeButtonPin, fakeLedPin)
	go bm.Run()
	// Check that a single 'pressed' event was sent
	select {
	case <-bm.GetChannel():
		return
	case <-time.After(time.Second * 5):
		t.Fatal("No pressed event detected")
	}
}
