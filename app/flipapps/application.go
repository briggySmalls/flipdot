package flipapps

import (
	fmt "fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/briggySmalls/flipdot/app/text"
	"golang.org/x/image/font"
)

const (
	messageInSize = 20
)

type application struct {
	flipdot       flipdot.Flipdot
	buttonManager ButtonManager
	imager        Imager
	// Externally-visible channel for adding messages to the application
	MessagesIn chan MessageRequest
}

// Creates and initialises a new Application
func NewApplication(flipdot flipdot.Flipdot, buttonManager ButtonManager, tickPeriod time.Duration, font font.Face) application {
	// Create a text builder
	width, height := flipdot.Size()
	textBuilder := text.NewTextBuilder(width, height, font)
	// Create the status image
	statusImage, err := readStatusImage("./bell.png")
	errorHandler(err)
	app := application{
		flipdot:       flipdot,
		buttonManager: buttonManager,
		imager:        NewImager(textBuilder, statusImage, uint(len(flipdot.Signs()))),
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
				images, err := s.imager.Clock(t, true)
				errorHandler(err)
				err = s.flipdot.Draw(images)
				errorHandler(err)
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
		// Create images from message
		images, err := s.imager.Message(message.From, message.GetText())
		errorHandler(err)
		// Send images
		s.sendImages(images)
	default:
		err = fmt.Errorf("Neither images or text supplied")
	}
	// Handle errors
	errorHandler(err)
}

// Helper function to send images to the signs
func (s *application) sendImages(images []*flipdot.Image) (err error) {
	err = s.flipdot.Draw(images)
	return
}

func readStatusImage(filename string) (image image.Image, err error) {
	// Read the image from disk
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	// Interpret as an image
	image, err = png.Decode(file)
	if err != nil {
		return
	}
	return
}

// In-queue error handler
func errorHandler(err error) {
	if err != nil {
		fmt.Printf("Runtime error: %s", err)
	}
}
