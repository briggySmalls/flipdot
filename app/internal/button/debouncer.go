package button

type DebounceState uint8

const (
	Off DebounceState = iota
	Transitioning
	Triggered
)

const (
	debouncedCount = 50
)

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
