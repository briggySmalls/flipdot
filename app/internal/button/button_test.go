package button

import (
	"testing"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"

	gomock "github.com/golang/mock/gomock"
)

func TestActive(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	// Create a button manager
	flashPeriod := time.Millisecond
	bm := NewButtonManager(fakeButtonPin, fakeLedPin, flashPeriod, time.Hour)

	// Configure mocks
	callCount := 0
	gomock.InOrder(
		fakeLedPin.EXPECT().Low(),
		fakeLedPin.EXPECT().Toggle().AnyTimes().Do(func() {
			callCount++
		}),
	)
	// fakeButtonPin.EXPECT().Read().AnyTimes() // Tested properly elsewhere
	// Simulate a call to turn on the LED
	bm.SetState(Active)
	multiple := 10
	time.Sleep(flashPeriod * time.Duration(multiple))
	// Wait for results
	if callCount < multiple-1 || callCount > multiple+1 {
		t.Errorf("Unexpected Toggle() call count %d", callCount)
	}
}

func TestButtonPressed(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	debounceDuration := time.Millisecond * 50
	// Configure mock to expect calls
	fakeLedPin.EXPECT().Low()
	fakeLedPin.EXPECT().Toggle().AnyTimes() // We test this properly in another test
	start := time.Now()
	fakeButtonPin.EXPECT().Read().AnyTimes().DoAndReturn(func() rpio.State {
		runtime := time.Now().Sub(start)
		if runtime > time.Millisecond && runtime < time.Millisecond+debounceDuration*2 {
			return rpio.High
		}
		return rpio.Low
	})
	// Create a button manager and start the app
	bm := NewButtonManager(fakeButtonPin, fakeLedPin, time.Hour, time.Millisecond*50)
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

func TestDeadlock(t *testing.T) {
	// Create fake buttons
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeLedPin := NewMockOutputPin(ctrl)
	fakeButtonPin := NewMockTriggerPin(ctrl)
	debounceDuration := time.Duration(time.Millisecond * 50)
	// Configure mock to expect calls
	fakeLedPin.EXPECT().Low().AnyTimes()
	fakeLedPin.EXPECT().Toggle().AnyTimes()
	fakeButtonPin.EXPECT().Read().AnyTimes().Return(rpio.High)
	// Create a button manager and start the app
	bm := NewButtonManager(fakeButtonPin, fakeLedPin, time.Hour, debounceDuration)
	bm.SetState(Active)
	// Wait to ensure that button is pressed
	time.Sleep(debounceDuration * 2)
	// Send state change request, this should not deadlock
	bm.SetState(Inactive)
}
