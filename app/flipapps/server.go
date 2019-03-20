package flipapps

import (
	context "context"
	fmt "fmt"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/status"
)

const (
	tokenDuration = time.Hour // Duration before JWT expiry
)

func NewRpcServer(secret, password string, messageQueue chan MessageRequest, signsInfo []*flipdot.GetInfoResponse_SignInfo) (grpcServer *grpc.Server) {
	// Create a flipdot server
	server := NewServer(secret, password, messageQueue, signsInfo)
	// create a gRPC server object
	grpcServer = grpc.NewServer(grpc.UnaryInterceptor(server.(*flipappsServer).unaryAuthInterceptor))
	// attach the FlipApps service to the server
	RegisterFlipAppsServer(grpcServer, server)
	return grpcServer
}

// Create a new server
func NewServer(secret, password string, messageQueue chan MessageRequest, signsInfo []*flipdot.GetInfoResponse_SignInfo) FlipAppsServer {
	// Create a flipdot controller
	server := &flipappsServer{
		appSecret:    secret,
		appPassword:  password,
		messageQueue: messageQueue,
		signsInfo:    signsInfo,
	}
	// Return the server
	return server
}

type flipappsServer struct {
	appSecret   string
	appPassword string
	// Channel to which new messages are sent
	messageQueue chan MessageRequest
	// Information on connected signs
	signsInfo []*flipdot.GetInfoResponse_SignInfo
}

// Handler for client request to authenticate (obtain JWT token)
func (f *flipappsServer) Authenticate(_ context.Context, request *AuthenticateRequest) (*AuthenticateResponse, error) {
	// Confirm the password is correct
	if request.Password != f.appPassword {
		return nil, status.Error(codes.Unauthenticated, "Incorrect password")
	}
	// Create a new token object, specifying signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(tokenDuration).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(f.appSecret))
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create authentication token")
	}
	return &AuthenticateResponse{Token: tokenString}, nil
}

// Handler for client request of information on connected signs
func (f *flipappsServer) GetInfo(_ context.Context, _ *flipdot.GetInfoRequest) (*flipdot.GetInfoResponse, error) {
	// Make a request to the controller
	signs := f.signsInfo
	response := flipdot.GetInfoResponse{Signs: signs}
	return &response, nil
}

// Handler for client request to display a message
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

// Interceptor that checks all RPC calls are authorized
func (f *flipappsServer) unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// First, check this isn't an auth call itself
	if info.FullMethod == fmt.Sprintf("/%s/%s", _FlipApps_serviceDesc.ServiceName, "Authenticate") {
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
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Unexpected signing method: %v", token.Header["alg"])
		}
		// Return secret key for parsing with
		return []byte(f.appSecret), nil
	})
	// Indicate if we are happy with the result
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Could not parse token: %s", t)
	}
	// Check claims are valid
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return status.Error(codes.Unauthenticated, "Invalid/expired token")
	}
	return nil
}
