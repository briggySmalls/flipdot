package flipapps

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/briggySmalls/flipcli/flipdot"
)

func NewFlipappsServer(flipdot flipdot.Flipdot) FlipAppsServer {
	// Create a flipdot controller
	return &flipappsServer{flipdot: flipdot}
}

type flipappsServer struct {
	flipdot flipdot.Flipdot
}

func (f *flipappsServer) GetInfo(_ context.Context, _ *flipdot.GetInfoRequest) (*flipdot.GetInfoResponse, error) {
	// Make a request to the controller
	err := f.flipdot.GetInfo()
	grpcError := status.Error(codes.Internal, err.Error())
	response := flipdot.GetInfoResponse{}
	return
}

func (f *flipappsServer) Draw(ctx context.Context, request *flipdot.DrawRequest) (response *flipdot.DrawResponse, err error) {
	// Make a request to the controller
	response, err = f.flipdot.Draw(request.Image)
	return
}

func (f *flipappsServer) Text(ctx context.Context, request *TextRequest) (response *TextResponse, err error) {
	err = f.flipdot.Text(request.Text)
}
