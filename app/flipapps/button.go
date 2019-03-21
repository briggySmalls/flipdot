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
	On
)

const (
	debouncedCount = 10
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

func NewTriggerPin(pinNum uint8, edge rpio.Edge) TriggerPin {
	// Create the pin
	p := rpio.Pin(pinNum)
	// Configure the pin
	p.Input()
	p.PullDown()
	p.Detect(edge)
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
	var debounceCounter uint16
	debounceState := Off
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
				pinState := b.buttonPin.Read()
				if debounceState == Off && pinState == rpio.High {
					// Start counter (we've transitioned)
					debounceState = Transitioning
					debounceCounter = 0
				} else if debounceState == Transitioning {
					if pinState == rpio.High {
						// Increment counter
						debounceCounter++
						if debounceCounter == debouncedCount {
							log.Println("Button press detected")
							b.buttonPressed <- struct{}{}
						}
					} else {
						// Reset state (press aborted)
						debounceState = Off
					}
				}
			}
		}
	}
}

func (b *buttonManager) GetChannel() chan struct{} {
	return b.buttonPressed
}
