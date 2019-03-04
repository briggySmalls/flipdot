package flipapps

import (
	"google.golang.org/grpc/codes"
	context "context"

	"github.com/briggySmalls/flipcli/flipdot"
)

func NewFlipappsServer(flipdot flipdot.Flipdot) FlipAppsServer {
	// Create a flipdot controller
	return flipappsServer{
		flipdot: flipdot
	}
}

type flipappsServer struct {
	flipdot flipdot.Flipdot
}

func (f *flipappsServer) GetInfo(_ context.Context, _ *flipdot.GetInfoRequest) (*GetInfoResponse, error) {
	// Make a request to the controller
	err := f.flipdot.GetInfo()
	grpcError := status.Error(codes.Internal, err.Error())
	response := GetInfoResponse{}
	return
}

func (f *flipappsServer) Draw(ctx context.Context, request *flipdot.DrawRequest) (response *DrawResponse, err error) {
	// Make a request to the controller
	err = f.flipdot.Draw(request.Image)
	return
}

func (f *flipappsServer) Text(ctx context.Context, request *TextRequest) (response *TextResponse, err error) {
	err = f.flipdot.Text(request.Text)
}
