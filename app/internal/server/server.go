package server

import (
	context "context"
	fmt "fmt"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/briggySmalls/flipdot/app/internal/protos"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/status"
)

func NewRpcServer(secret, password string, tokenExpiry time.Duration, messageQueue chan protos.MessageRequest, signsInfo []*protos.GetInfoResponse_SignInfo) (grpcServer *grpc.Server) {
	// Create a flipdot server
	server := NewServer(secret, password, tokenExpiry, messageQueue, signsInfo)
	// create a gRPC server object
	grpcServer = grpc.NewServer(grpc.UnaryInterceptor(server.(*flipappsServer).unaryAuthInterceptor))
	// attach the FlipApps service to the server
	protos.RegisterFlipAppsServer(grpcServer, server)
	return grpcServer
}

// Create a new server
func NewServer(secret, password string, tokenExpiry time.Duration, messageQueue chan protos.MessageRequest, signsInfo []*protos.GetInfoResponse_SignInfo) protos.FlipAppsServer {
	// Create a flipdot controller
	server := &flipappsServer{
		appSecret:    secret,
		appPassword:  password,
		tokenExpiry:  tokenExpiry,
		messageQueue: messageQueue,
		signsInfo:    signsInfo,
	}
	// Return the server
	return server
}

type flipappsServer struct {
	appSecret   string
	appPassword string
	// Time after which an authorisation token expires
	tokenExpiry time.Duration
	// Channel to which new messages are sent
	messageQueue chan protos.MessageRequest
	// Information on connected signs
	signsInfo []*protos.GetInfoResponse_SignInfo
}

// Handler for client request to authenticate (obtain JWT token)
func (f *flipappsServer) Authenticate(_ context.Context, request *protos.AuthenticateRequest) (*protos.AuthenticateResponse, error) {
	// Confirm the password is correct
	if request.Password != f.appPassword {
		return nil, status.Error(codes.Unauthenticated, "Incorrect password")
	}
	// Create a new token object, specifying signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(f.tokenExpiry).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(f.appSecret))
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create authentication token")
	}
	return &protos.AuthenticateResponse{Token: tokenString}, nil
}

// Handler for client request of information on connected signs
func (f *flipappsServer) GetInfo(_ context.Context, _ *protos.GetInfoRequest) (*protos.GetInfoResponse, error) {
	// Make a request to the controller
	signs := f.signsInfo
	response := protos.GetInfoResponse{Signs: signs}
	return &response, nil
}

// Handler for client request to display a message
func (f *flipappsServer) SendMessage(ctx context.Context, request *protos.MessageRequest) (response *protos.MessageResponse, err error) {
	switch request.Payload.(type) {
	case *protos.MessageRequest_Images, *protos.MessageRequest_Text:
		// Enqueue message
		f.messageQueue <- *request
	default:
		err = status.Error(codes.InvalidArgument, "Neither images or text supplied")
	}
	response = &protos.MessageResponse{}
	return
}

// Interceptor that checks all RPC calls are authorized
func (f *flipappsServer) unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// First, check this isn't an auth call itself
	if info.FullMethod == fmt.Sprintf("/flipapps.FlipApps/Authenticate") {
		// We don't need to check for tokens here
		return handler(ctx, req)
	}
	// Try to pull out token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// Caller didn't supply a token
		return nil, status.Error(codes.Unauthenticated, "Authentication token not provided")
	}
	if len(md["token"]) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Badly formatted metadata (missing token)")
	}
	// Check the token
	err := f.checkToken(md["token"][0])
	if err != nil {
		return nil, err
	}
	// Execute the usual RPC clal
	return handler(ctx, req)
}

// Helper function to check a request's JWT token is valid
func (f *flipappsServer) checkToken(t string) error {
	// Parse JWT token
	_, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Unexpected signing method: %v", token.Header["alg"])
		}
		// Return secret key for parsing with
		return []byte(f.appSecret), nil
	})
	// Indicate if we are happy with the result
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "%s", err)
	}
	return nil
}
