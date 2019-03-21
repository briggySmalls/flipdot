package flipapps

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type State uint8
type DebounceState uint8

const (
	Active State = iota
	Inactive
	Stopped
)

const (
	Off DebounceState = iota
	Transitioning
	Triggered
)

const (
	debouncedCount = 50
)

type TriggerPin interface {
	EdgeDetected() bool
	Read() rpio.State
}

type OutputPin interface {
	High()
	Low()
	Toggle()
}

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

func NewTriggerPin(pinNum uint8) TriggerPin {
	// Create the pin
	p := rpio.Pin(pinNum)
	// Configure the pin
	p.Input()
	p.PullDown()
	return p
}

func NewOutputPin(pinNum uint8) OutputPin {
	// Create the pin
	p := rpio.Pin(pinNum)
	// Configure the pin
	p.Output()
	return p
}

func NewButtonManager(buttonPin TriggerPin, ledPin OutputPin, flashFreq, debounceFreq time.Duration) ButtonManager {
	manager := buttonManager{
		buttonPressed: make(chan struct{}),
		stateChanger:  make(chan State),
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
		// Update manager state (thread safe?)
		b.stateChanger <- state
	}
}

func (b *buttonManager) run(flashFreq time.Duration, debounceTime time.Duration) {
	// Run control loop
	log.Println("Button manager loop starting...")
	b.ledPin.Low() // Ensure LED off
	flashTicker := time.NewTicker(flashFreq)

	// Debounce-releated stuff
	debounceTicker := time.NewTicker(debounceTime / debouncedCount)
	debouncer := debouncer{state: Off, debouncedCount: debouncedCount}
	for {
		select {
		case state, ok := <-b.stateChanger:
			if !ok {
				// Stopping, tell listeners by stopping 'pressed' channel
				close(b.buttonPressed)
				return
			}
			b.state = state
			log.Printf("State updated to %d", b.state)
			if b.state == Inactive {
				// Ensure LED is not illuminated
				b.ledPin.Low()
			}
		case <-flashTicker.C:
			if b.state == Active {
				// Toggle LED illumination on tick
				b.ledPin.Toggle()
			}
		case <-debounceTicker.C:
			// Check if button status has changed
			if b.state == Active {
				// Read the button state
				pinState := b.buttonPin.Read()
				// Pass it to the debouncer
				if debouncer.debounce(pinState == rpio.High) {
					log.Println("Button press detected")
					b.buttonPressed <- struct{}{}
				}
			}
		}
	}
}

type debouncer struct {
	// Internal state of the debouncer
	state DebounceState
	// Counter for the number of consecutive positive signal counts
	counter uint16
	// Number of consecutive counts considered to be 'debounced'
	debouncedCount uint16
}

// Debounces an input boolan, returning the debounced signal
func (d *debouncer) debounce(value bool) bool {
	switch d.state {
	case Off:
		if value {
			// Start counter (we've transitioned)
			d.state = Transitioning
			d.counter = 0
		}
	case Transitioning:
		if value {
			// Increment counter
			d.counter++
			if d.counter == d.debouncedCount {
				// Sufficient consecutive high signals
				d.state = Triggered
				return true
			}
		} else {
			// Reset state (press aborted)
			d.state = Off
		}
	case Triggered:
		if !value {
			// Button has returned to off
			d.state = Off
		}
	}
	return false
}

func (b *buttonManager) GetChannel() chan struct{} {
	return b.buttonPressed
}
