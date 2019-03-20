package flipapps

import "github.com/stianeikeland/go-rpio"

type TriggerPin interface {
	EdgeDetected() bool
}

type OutputPin interface {
	High()
	Low()
}

type buttonManager struct {
	buttonPin TriggerPin
	ledPin    OutputPin
	// Channel read by this manager to set LED state
	buttonPressed chan struct{}
	// Internal channel to stop polling for edge detection
	done chan struct{}
}

type ButtonManager interface {
	Run()
	Stop()
	WriteLed(bool)
	GetChannel() chan struct{}
}

func NewTriggerPin(pinNum uint8, edge rpio.Edge) TriggerPin {
	// Create the pin
	p := rpio.Pin(pinNum)
	// Configure the pin
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

func NewButtonManager(buttonPin TriggerPin, ledPin OutputPin) ButtonManager {
	return &buttonManager{
		buttonPressed: make(chan struct{}),
		done:          make(chan struct{}),
		buttonPin:     buttonPin,
		ledPin:        ledPin,
	}
}

func (b *buttonManager) Run() {
	// Run control loop
	select {
	case <-b.done:
		// We've been asked to stop, shut 'pressed' channel
		close(b.buttonPressed)
		return
	default:
		// Check if button status has changed
		if b.buttonPin.EdgeDetected() {
			b.buttonPressed <- struct{}{}
		}
	}
}

func (b *buttonManager) Stop() {
	b.done <- struct{}{}
}

func (b *buttonManager) WriteLed(status bool) {
	if status {
		b.ledPin.High()
	} else {
		b.ledPin.Low()
	}
}

func (b *buttonManager) GetChannel() chan struct{} {
	return b.buttonPressed
}
