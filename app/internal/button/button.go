package button

import (
	"fmt"
	"log"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

type State uint8

const (
	Active State = iota
	Inactive
	Stopped
)

const (
	stateQueueCount = 10
)

type buttonManager struct {
	buttonPin TriggerPin
	ledPin    OutputPin
	// Channel read by this manager to set LED state
	buttonPressed chan struct{}
	// Internal channel to change state
	stateChanger chan State
	// Internal record of state
	state State
}

type ButtonManager interface {
	SetState(State)
	GetChannel() chan struct{}
}

func NewButtonManager(buttonPin TriggerPin, ledPin OutputPin, flashFreq, debounceFreq time.Duration) ButtonManager {
	manager := buttonManager{
		buttonPressed: make(chan struct{}),
		stateChanger:  make(chan State, stateQueueCount),
		buttonPin:     buttonPin,
		ledPin:        ledPin,
		state:         Inactive,
	}
	// Start the manager control loop
	go manager.run(flashFreq, debounceFreq)
	// Return the manager
	return &manager
}

func (b *buttonManager) SetState(state State) {
	if state == Stopped {
		// Stop the manager thread
		close(b.stateChanger)
	} else {
		// Update manager state
		// This shouldn't ever block, although in theory it could
		b.stateChanger <- state
	}
}

// Button manager blocking activity loop (designed to be run in goroutine)
func (b *buttonManager) run(flashFreq time.Duration, debounceTime time.Duration) {
	// Run control loop
	log.Println("Button manager loop starting...")
	b.ledPin.Low() // Ensure LED off

	// Create some channels for stopping 'active' goroutines
	stopButtonFlashing := make(chan struct{})
	stopButtonListening := make(chan struct{})
	defer close(stopButtonFlashing)
	defer close(stopButtonListening)

	// Keep popping state changes off the queue
	for {
		state, ok := <-b.stateChanger
		// Check if we need to stop
		if !ok {
			close(b.buttonPressed)
			return
		}
		// Check we're changing state
		if b.state == state {
			// We don't handle same-same transitions
			continue
		}
		// Record state change
		b.state = state
		// Handle state change
		switch state {
		// Button becomes active
		case Active:
			// Run goroutine for flashing button
			go togglePinPeriodically(b.ledPin, flashFreq, stopButtonFlashing)
			// Run goroutine for listening for a button press
			go listenForPress(b.buttonPin, debounceTime, stopButtonListening, b.buttonPressed)
		// Button becomes inactive
		case Inactive:
			// Stop goroutines
			stopButtonListening <- struct{}{}
			stopButtonFlashing <- struct{}{}
		// Unexpected state value
		default:
			panic(fmt.Errorf("Unexpected state"))
		}
	}
}

// Function for toggling a pin until stopped
func togglePinPeriodically(pin OutputPin, flashFreq time.Duration, done <-chan struct{}) {
	// Create a ticker
	ticker := time.NewTicker(flashFreq)
	defer ticker.Stop()
	// Toggle until we have to stop
	for {
		select {
		// Toggle the pin when the ticker elapses
		case <-ticker.C:
			pin.Toggle()
		// Set pin low and terminate when told we are done
		case <-done:
			pin.Low()
			return
		}
	}
}

// Listen for a button press
func listenForPress(pin TriggerPin, debounceTime time.Duration, done <-chan struct{}, pressed chan<- struct{}) {
	// Create a ticker to poll button state
	debounceTicker := time.NewTicker(debounceTime / debouncedCount)
	defer debounceTicker.Stop()
	// Create debouncer that keeps track of previous activity
	debouncer := debouncer{state: Off, debouncedCount: debouncedCount}
	// Run
	for {
		select {
		case <-debounceTicker.C:
			// Read the button state
			pinState := pin.Read()
			// Pass it to the debouncer
			if debouncer.debounce(pinState == rpio.High) {
				log.Println("Button press detected")
				pressed <- struct{}{}
			}
		case <-done:
			// We've been told we're done
			return
		}
	}
}

func (b *buttonManager) GetChannel() chan struct{} {
	return b.buttonPressed
}
