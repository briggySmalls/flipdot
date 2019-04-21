package server

import (
	context "context"
	reflect "reflect"
	"testing"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	contextTimeoutS = 1
	password        = "password"
)

func TestAuthenticateFail(t *testing.T) {
	flipapps, queue, _ := createTestObjects(t)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	_, err := flipapps.Authenticate(ctx, &protos.AuthenticateRequest{Password: "wrong"})
	// Check respose
	if err == nil {
		t.Fatal("Failed to detect failed password")
	}
	if s, ok := status.FromError(err); !ok || s.Code() != codes.Unauthenticated {
		t.Errorf("Failed to assign appropriate error code: %s", s.Code())
	}
	// Check no messages were sent
	checkNoMessages(t, queue)
}

func TestAuthenticatePass(t *testing.T) {
	flipapps, queue, _ := createTestObjects(t)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	response, err := flipapps.Authenticate(ctx, &protos.AuthenticateRequest{Password: password})
	// Check response
	if err != nil {
		t.Fatal(err)
	}
	// Assert we got a token back
	if response.Token == "" {
		t.Error("Failed to return token")
	}
	// Assert we can roudtrip the token
	if flipapps.(*flipappsServer).checkToken(response.Token) != nil {
		t.Error("Failed to check token")
	}
	// Check no messages were sent
	checkNoMessages(t, queue)
}

func TestGetInfo(t *testing.T) {
	flipapps, queue, signs := createTestObjects(t)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	response, err := flipapps.GetInfo(ctx, &protos.GetInfoRequest{})
	// Assert the return values
	if err != nil {
		t.Fatal(err)
	}
	// Ensure every sign matches exactly
	if !reflect.DeepEqual(response.Signs, signs) {
		t.Errorf("Signs don't match")
	}
	// Check no messages were sent
	checkNoMessages(t, queue)
}

func TestSendMessage(t *testing.T) {
	// Create mocks
	flipapps, queue, _ := createTestObjects(t)
	// Run the command
	ctx, cancel := getContext()
	defer cancel()
	originalRequest := protos.MessageRequest{
		From:    "briggySmalls",
		Payload: &protos.MessageRequest_Text{"test text"},
	}
	_, err := flipapps.SendMessage(ctx, &originalRequest)
	// Assert the return values
	if err != nil {
		t.Fatal(err)
	}
	// Confirm that the message was enqueued
	select {
	case message := <-queue:
		// Assert message is as expected
		reflect.DeepEqual(originalRequest, message)
	default:
		// No message enqueued
		t.Fatal("No message was enqueued")
	}
}

// Helper function to set up the unit under test
func createTestObjects(t *testing.T) (protos.FlipAppsServer, chan protos.MessageRequest, []*protos.GetInfoResponse_SignInfo) {
	// Make some dummy signs
	signs := []*protos.GetInfoResponse_SignInfo{
		{
			Name: "test1", Width: 10, Height: 2,
		},
		{
			Name: "test2", Width: 10, Height: 2,
		},
	}
	// Make a channel for sending messages
	messageQueue := make(chan protos.MessageRequest, 10)
	// Create object under test
	server := NewServer("secret", "password", time.Hour, messageQueue, signs)
	return server, messageQueue, signs
}

// Helper function to check that no messages were queued by the server
func checkNoMessages(t *testing.T, queue chan protos.MessageRequest) {
	// Check no messages were sent
	select {
	case message := <-queue:
		t.Fatalf("Unexpected message received from %s", message.From)
	default:
		// We expect there to be no messages
		return
	}
}

// Get a simple context for sending requests via gRPC
func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), contextTimeoutS*time.Second)
}
