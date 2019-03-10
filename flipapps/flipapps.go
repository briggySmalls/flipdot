package flipapps

import (
	context "context"
	fmt "fmt"

	"google.golang.org/grpc/codes"

	"github.com/briggySmalls/flipapp/flipdot"
	"golang.org/x/image/font"
	"google.golang.org/grpc/status"
)

const (
	messageQueueSize = 20
)

// Create a new server
func NewFlipappsServer(flipdot flipdot.Flipdot, font font.Face) FlipAppsServer {
	// Create a flipdot controller
	server := &flipappsServer{
		flipdot:      flipdot,
		font:         font,
		messageQueue: make(chan MessageRequest, messageQueueSize),
	}
	// Run the queue pump
	go server.run()
	// Return the server
	return server
}

type flipappsServer struct {
	flipdot      flipdot.Flipdot
	font         font.Face
	messageQueue chan MessageRequest
}

// Get info about connected signs
func (f *flipappsServer) GetInfo(_ context.Context, _ *flipdot.GetInfoRequest) (*flipdot.GetInfoResponse, error) {
	// Make a request to the controller
	signs := f.flipdot.Signs()
	response := flipdot.GetInfoResponse{Signs: signs}
	return &response, nil
}

// Send a message to be displayed on the signs
func (f *flipappsServer) SendMessage(ctx context.Context, request *MessageRequest) (response *MessageResponse, err error) {
	switch request.Payload.(type) {
	case *MessageRequest_Images, *MessageRequest_Text:
		f.messageQueue <- *request
	default:
		err = status.Error(codes.InvalidArgument, "Neither images or text supplied")
	}
	response = &MessageResponse{}
	return
}

// Helper function to send text to the signs
func (f *flipappsServer) sendText(txt string) (err error) {
	err = f.flipdot.Text(txt, f.font, false)
	return
}

// Helper function to send images to the signs
func (f *flipappsServer) sendImages(images []*flipdot.Image) (err error) {
	err = f.flipdot.Draw(images)
	return
}

// Routine for handling queued messages
func (f *flipappsServer) run() {
	// Run forever
	for {
		// Pop a message off the queue
		message := <-f.messageQueue
		// Handle message depending on type (text/images)
		var err error
		switch message.Payload.(type) {
		case *MessageRequest_Images:
			err = f.sendImages(message.GetImages().Images)
		case *MessageRequest_Text:
			err = f.sendText(message.GetText())
		default:
			err = status.Error(codes.InvalidArgument, "Neither images or text supplied")
		}
		// Handle errors
		errorHandler(err)
	}
}

// In-queue error handler
func errorHandler(err error) {
	if err != nil {
		fmt.Printf("Runtime error: %s", err)
	}
}
