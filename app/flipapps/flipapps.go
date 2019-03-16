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
	"golang.org/x/image/font"
	"google.golang.org/grpc/status"
)

const (
	messageQueueSize = 20
	tokenDuration    = time.Hour
)

func NewRpcFlipappsServer(flipdot flipdot.Flipdot, font font.Face, secret, password string) (grpcServer *grpc.Server) {
	// Create a flipdot server
	server := NewFlipappsServer(flipdot, font, secret, password)
	// create a gRPC server object
	grpcServer = grpc.NewServer(grpc.UnaryInterceptor(server.(*flipappsServer).unaryAuthInterceptor))
	// attach the FlipApps service to the server
	RegisterFlipAppsServer(grpcServer, server)
	return grpcServer
}

// Create a new server
func NewFlipappsServer(flipdot flipdot.Flipdot, font font.Face, secret, password string) FlipAppsServer {
	// Create a flipdot controller
	server := &flipappsServer{
		flipdot:      flipdot,
		font:         font,
		messageQueue: make(chan MessageRequest, messageQueueSize),
		appSecret:    secret,
		appPassword:  password,
	}
	// Run the queue pump concurrently
	go server.run()
	// Return the server
	return server
}

type flipappsServer struct {
	flipdot      flipdot.Flipdot
	font         font.Face
	messageQueue chan MessageRequest
	appSecret    string
	appPassword  string
}

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

// Routine for handling queued messages
func (f *flipappsServer) run() {
	// Create a ticker
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	pause := false
	// Run forever
	for {
		select {
		// Handle message, if available
		case message := <-f.messageQueue:
			// Pause clock whilst we handle a message
			pause = true
			// Handle message
			f.handleMessage(message)
			// Unpause clock
			pause = false
		// Otherwise display the time
		case t := <-ticker.C:
			// Only display the time if we've not paused the clock
			if !pause {
				// Print the time (centred)
				f.sendText(t.Format("Mon 1 Jan\n3:04 PM"), true)
			}
		}
	}
}

// Check that all RPC calls are authorized
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
	// Parse JWT token
	token, err := jwt.Parse(md["token"][0], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.InvalidArgument, "Unexpected signing method: %v", token.Header["alg"])
		}
		// Return secret key for parsing with
		return f.appSecret, nil
	})
	// Indicate if we are happy with the result
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Could not parse token")
	}
	// Check claims are valid
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "Invalid/expired token")
	}
	// Execute the usual RPC clal
	return handler(ctx, req)
}

func (f *flipappsServer) handleMessage(message MessageRequest) {
	var err error
	switch message.Payload.(type) {
	case *MessageRequest_Images:
		err = f.sendImages(message.GetImages().Images)
	case *MessageRequest_Text:
		// Send left-aligned text for messages
		err = f.sendText(message.GetText(), false)
	default:
		err = status.Error(codes.InvalidArgument, "Neither images or text supplied")
	}
	// Handle errors
	errorHandler(err)
}

// Helper function to send text to the signs
func (f *flipappsServer) sendText(txt string, center bool) (err error) {
	err = f.flipdot.Text(txt, f.font, center)
	return
}

// Helper function to send images to the signs
func (f *flipappsServer) sendImages(images []*flipdot.Image) (err error) {
	err = f.flipdot.Draw(images)
	return
}

// In-queue error handler
func errorHandler(err error) {
	if err != nil {
		fmt.Printf("Runtime error: %s", err)
	}
}
