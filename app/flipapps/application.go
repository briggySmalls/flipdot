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
	messageQueueSize = 20
)

type application struct {
	flipdot      flipdot.Flipdot
	tickPeriod   time.Duration
	font         font.Face
	MessageQueue chan MessageRequest
}

// Creates and initialises a new Application
func NewApplication(flipdot flipdot.Flipdot, tickPeriod time.Duration, font font.Face) application {
	return application{
		flipdot:      flipdot,
		tickPeriod:   tickPeriod,
		font:         font,
		MessageQueue: make(chan MessageRequest, messageQueueSize),
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
		// Handle message, if available
		case message, open := <-s.MessageQueue:
			if !open {
				// No more messages to handle
				return
			}
			// Pause clock whilst we handle a message
			pause = true
			// Handle message
			s.handleMessage(message)
			// Unpause clock
			pause = false
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
