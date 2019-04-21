package button

import rpio "github.com/stianeikeland/go-rpio/v4"

type TriggerPin interface {
	Read() rpio.State
}

type OutputPin interface {
	High()
	Low()
	Toggle()
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
