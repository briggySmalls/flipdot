package flipapps

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type State uint8

const (
	Active State = iota
	Inactive
	Stopped
)

type TriggerPin interface {
	EdgeDetected() bool
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
	state     State
	frequency time.Duration
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

func NewButtonManager(buttonPin TriggerPin, ledPin OutputPin, frequency time.Duration) ButtonManager {
	manager := buttonManager{
		buttonPressed: make(chan struct{}),
		stateChanger:  make(chan State),
		buttonPin:     buttonPin,
		ledPin:        ledPin,
		state:         Inactive,
		frequency:     frequency,
	}
	// Start the manager control loop
	go manager.run()
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

func (b *buttonManager) run() {
	// Run control loop
	log.Println("Button manager loop starting...")
	b.ledPin.Low() // Ensure LED off
	ticker := time.NewTicker(b.frequency)
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
			// Ensure LED is not illuminated
			b.ledPin.Low()
		case <-ticker.C:
			if b.state == Active {
				// Toggle LED illumination on tick
				b.ledPin.Toggle()
			}
		default:
			// Check if button status has changed
			if b.state == Active && b.buttonPin.EdgeDetected() {
				log.Println("Button press detected")
				b.buttonPressed <- struct{}{}
			}
		}
	}
}

func (b *buttonManager) GetChannel() chan struct{} {
	return b.buttonPressed
}
