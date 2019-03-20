package flipapps

import (
	fmt "fmt"
	"time"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"golang.org/x/image/font"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	messageInSize = 20
)

type application struct {
	flipdot       flipdot.Flipdot
	buttonManager ButtonManager
	tickPeriod    time.Duration
	font          font.Face
	// Externally-visible channel for adding messages to the application
	MessagesIn chan MessageRequest
	// Staging channel used by application internally
	messagesPending chan MessageRequest
	ShowMessage     chan struct{}
}

// Creates and initialises a new Application
func NewApplication(flipdot flipdot.Flipdot, buttonManager ButtonManager, tickPeriod time.Duration, font font.Face) application {
	return application{
		flipdot:         flipdot,
		buttonManager:   buttonManager,
		tickPeriod:      tickPeriod,
		font:            font,
		MessagesIn:      make(chan MessageRequest, messageInSize),
		messagesPending: make(chan MessageRequest, messageInSize),
		ShowMessage:     make(chan struct{}),
	}
}

// Routine for handling queued messages
func (s *application) Run() {
	// Create a ticker
	ticker := time.NewTicker(s.tickPeriod)
	defer ticker.Stop()
	pause := false

	// Run forever
	for {
		select {
		case message, ok := <-s.MessagesIn:
			if !ok {
				// There will be no more messages to handle
				// TODO: check if there are any pending message that we should wait for
				return
			}
			// Externally queued message is available
			// Pass to internal buffer
			s.messagesPending <- message
			// We have at least one message, so light LED
			s.buttonManager.WriteLed(true)
		// Handle user signal to display message
		case _, ok := <-s.ShowMessage:
			if !ok {
				// There will be no more show-message requests
				return
			}
			// Check if there are pending messages
			select {
			case message, ok := <-s.messagesPending:
				if !ok {
					// There will be no more message requests
					return
				}
				// We have a message waiting, so display it
				pause = true // Pause clock whilst we handle a message
				// Handle message
				s.handleMessage(message)
				// Unpause clock
				pause = false
			default:
				// No message waiting, so skip
				// Also update the LED to indicate no more messages
				s.buttonManager.WriteLed(false)
				continue
			}
		// Otherwise display the time
		case t := <-ticker.C:
			// Only display the time if we've not paused the clock
			if !pause {
				// Print the time (centred)
				s.sendText(t.Format("Mon 1 Jan\n3:04 PM"), true)
			}
		}
	}
}

// Gets a message sent to the flipdot signs
func (s *application) handleMessage(message MessageRequest) {
	var err error
	switch message.Payload.(type) {
	case *MessageRequest_Images:
		err = s.sendImages(message.GetImages().Images)
	case *MessageRequest_Text:
		// Send left-aligned text for messages
		err = s.sendText(message.GetText(), false)
	default:
		err = status.Error(codes.InvalidArgument, "Neither images or text supplied")
	}
	// Handle errors
	errorHandler(err)
}

// Helper function to send text to the signs
func (s *application) sendText(txt string, center bool) (err error) {
	err = s.flipdot.Text(txt, s.font, center)
	return
}

// Helper function to send images to the signs
func (s *application) sendImages(images []*flipdot.Image) (err error) {
	err = s.flipdot.Draw(images)
	return
}

// In-queue error handler
func errorHandler(err error) {
	if err != nil {
		fmt.Printf("Runtime error: %s", err)
	}
}
