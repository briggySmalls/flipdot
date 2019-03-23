package flipapps

import (
	fmt "fmt"
	"log"
	"time"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"golang.org/x/image/font"
)

const (
	messageInSize = 20
)

type application struct {
	flipdot       flipdot.Flipdot
	buttonManager ButtonManager
	font          font.Face
	// Externally-visible channel for adding messages to the application
	MessagesIn chan MessageRequest
}

// Creates and initialises a new Application
func NewApplication(flipdot flipdot.Flipdot, buttonManager ButtonManager, tickPeriod time.Duration, font font.Face) application {
	app := application{
		flipdot:       flipdot,
		buttonManager: buttonManager,
		font:          font,
		MessagesIn:    make(chan MessageRequest, messageInSize),
	}
	go app.run(tickPeriod)
	return app
}

// Routine for handling queued messages
func (s *application) run(tickPeriod time.Duration) {
	// Create a ticker
	log.Println("Starting application loop...")
	pause := false
	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()
	// Create intermediate message queue
	pendingMessages := make([]MessageRequest, 0)
	// Get queue for button presses
	buttonPressed := s.buttonManager.GetChannel()
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
			log.Println("Message received")
			// Pass to internal buffer
			pendingMessages = append(pendingMessages, message)
			// We have at least one message, so activate button
			s.buttonManager.SetState(Active)
		// Handle user signal to display message
		case <-buttonPressed:
			log.Println("Show message request")
			// Check if there are pending messages
			if len(pendingMessages) > 0 {
				log.Println("Displaying message")
				// Disable button whilst we show a message
				s.buttonManager.SetState(Inactive)
				// Pop message
				message := pendingMessages[0]
				pendingMessages = pendingMessages[1:]
				// Display message
				s.handleMessage(message)
				// Reenable button if there are more messages
				if len(pendingMessages) > 0 {
					s.buttonManager.SetState(Active)
				}
			}
		// Otherwise display the time
		case t := <-ticker.C:
			// Only display the time if we've not paused the clock
			if !pause {
				log.Println("Tick event")
				// Print the time (centred)
				s.sendText(t.Format("Mon 1 Jan\n3:04 pm"), true)
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
		err = fmt.Errorf("Neither images or text supplied")
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
