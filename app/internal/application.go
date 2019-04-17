package internal

import (
	fmt "fmt"
	"log"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/button"
	"github.com/briggySmalls/flipdot/app/internal/client"
	"github.com/briggySmalls/flipdot/app/internal/imaging"
	"github.com/briggySmalls/flipdot/app/internal/server"
	"github.com/briggySmalls/flipdot/app/internal/shared"
)

const (
	messageInSize = 20
)

type application struct {
	flipdot       client.Flipdot
	buttonManager button.ButtonManager
	imager        imaging.Imager
	// Externally-visible channel for adding messages to the application
	messagesIn chan server.MessageRequest
}

type Application interface {
	GetMessagesChannel() chan server.MessageRequest
	Run(tickPeriod time.Duration)
}

// Creates and initialises a new Application
func NewApplication(flipdot client.Flipdot, buttonManager button.ButtonManager, imager imaging.Imager) Application {
	app := application{
		flipdot:       flipdot,
		buttonManager: buttonManager,
		imager:        imager,
		messagesIn:    make(chan server.MessageRequest, messageInSize),
	}
	return &app
}

func (a *application) GetMessagesChannel() chan server.MessageRequest {
	return a.messagesIn
}

// Blocking call that runs forever, polling for button presses, messages, and ticks
func (a *application) Run(tickPeriod time.Duration) {
	// Create a ticker
	log.Println("Starting application loop...")
	pause := false
	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()
	// Create intermediate message queue
	pendingMessages := make([]server.MessageRequest, 0)
	// Get queue for button presses
	buttonPressed := a.buttonManager.GetChannel()
	// Draw first clock
	a.drawTime(time.Now(), false)
	// Run forever
	for {
		select {
		case message, ok := <-a.messagesIn:
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
			a.buttonManager.SetState(button.Active)
			// Update time with message status
			a.drawTime(time.Now(), true)
		// Handle user signal to display message
		case <-buttonPressed:
			log.Println("Show message request")
			// Check if there are pending messages
			if len(pendingMessages) > 0 {
				log.Println("Displaying message")
				// Disable button whilst we show a message
				a.buttonManager.SetState(button.Inactive)
				// Pop message
				message := pendingMessages[0]
				pendingMessages = pendingMessages[1:]
				// Display message
				a.handleMessage(message)
				// Reenable button if there are more messages
				if len(pendingMessages) > 0 {
					a.buttonManager.SetState(button.Active)
				}
			}
		// Otherwise display the time
		case t := <-ticker.C:
			// Only display the time if we've not paused the clock
			if !pause {
				log.Println("Tick event")
				// Print the time (centred)
				a.drawTime(t, len(pendingMessages) > 0)
			}
		}
	}
}

// Helper function to draw the time on the signs
func (a *application) drawTime(time time.Time, isMessageAvailable bool) {
	// Print the time (centred)
	images, err := a.imager.Clock(time, isMessageAvailable)
	shared.ErrorHandler(err)
	err = a.flipdot.Draw(images, false)
	shared.ErrorHandler(err)
}

// Gets a message sent to the flipdot signs
func (a *application) handleMessage(message server.MessageRequest) {
	var err error
	switch message.Payload.(type) {
	case *server.MessageRequest_Images:
		err = a.sendImages(message.GetImages().Images)
	case *server.MessageRequest_Text:
		// Create images from message
		images, err := a.imager.Message(message.From, message.GetText())
		shared.ErrorHandler(err)
		// Send images
		a.sendImages(images)
	default:
		err = fmt.Errorf("Neither images or text supplied")
	}
	// Handle errors
	shared.ErrorHandler(err)
}

// Helper function to send images to the signs
func (a *application) sendImages(images []*client.Image) (err error) {
	err = a.flipdot.Draw(images, true)
	return
}
