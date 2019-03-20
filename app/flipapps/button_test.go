package flipapps

import (
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
)

func TestActive(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	// Create a button manager
	bm := NewButtonManager(fakeButtonPin, fakeLedPin, time.Microsecond)
	// Configure mock to expect calls
	done := make(chan struct{})
	fakeLedPin.EXPECT().Low().Times(2)
	fakeButtonPin.EXPECT().EdgeDetected().AnyTimes().Return(false)
	fakeLedPin.EXPECT().Toggle().Do(func() {
		done <- struct{}{}
	})
	// Simulate a call to turn on the LED
	bm.SetState(Active)
	// Wait for results
	select {
	case <-done:
		return
	case <-time.After(time.Second * 5):
		t.Fatal("No pressed event detected")
	}
}

func TestButtonPressed(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	// Configure mock to expect calls
	fakeLedPin.EXPECT().Low().Times(2)
	fakeButtonPin.EXPECT().EdgeDetected().AnyTimes().Return(false).Return(false).Return(true)
	fakeLedPin.EXPECT().Toggle().AnyTimes()
	// Create a button manager and start the app
	bm := NewButtonManager(fakeButtonPin, fakeLedPin, time.Microsecond)
	bm.SetState(Active)
	// Check that a single 'pressed' event was sent
	select {
	case <-bm.GetChannel():
		// Button pressed
		return
	case <-time.After(time.Second * 5):
		t.Fatal("No pressed event detected")
	}
}
