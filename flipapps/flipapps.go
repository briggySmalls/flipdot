package flipapps

import (
	context "context"

	"google.golang.org/grpc/codes"

	"github.com/briggySmalls/flipapp/flipdot"
	"golang.org/x/image/font"
	"google.golang.org/grpc/status"
)

func NewFlipappsServer(flipdot flipdot.Flipdot, font font.Face) FlipAppsServer {
	// Create a flipdot controller
	return &flipappsServer{
		flipdot: flipdot,
		font:    font,
	}
}

type flipappsServer struct {
	flipdot flipdot.Flipdot
	font    font.Face
}

func (f *flipappsServer) GetInfo(_ context.Context, _ *flipdot.GetInfoRequest) (*flipdot.GetInfoResponse, error) {
	// Make a request to the controller
	signs := f.flipdot.Signs()
	response := flipdot.GetInfoResponse{Signs: signs}
	return &response, nil
}

func (f *flipappsServer) SendMessage(ctx context.Context, request *MessageRequest) (response *MessageResponse, err error) {
	// Check if we are sending text or images
	switch request.Payload.(type) {
	case *MessageRequest_Images:
		err = f.sendImages(request.GetImages().Images)
	case *MessageRequest_Text:
		err = f.sendText(request.GetText())
	default:
		err = status.Error(codes.InvalidArgument, "Neither images or text supplied")
	}
	response = &MessageResponse{}
	return
}

func (f *flipappsServer) sendText(txt string) (err error) {
	err = f.flipdot.Text(txt, f.font, false)
	return
}

func (f *flipappsServer) sendImages(images []*flipdot.Image) (err error) {
	err = f.flipdot.Draw(images)
	return
}
