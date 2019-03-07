package main

import (
	"flag"
	fmt "fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"

	"github.com/briggySmalls/flipcli/flipapps"
	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/briggySmalls/flipcli/text"
	"golang.org/x/image/font"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Create a very simple command line
	clientPort := flag.Uint("client-port", 5001, "Port to connect to client on")
	serverPort := flag.Uint("server-port", 5002, "Port to expose server on")
	fontFile := flag.String("font", "", ".ttf font file for creating text images")
	fontSize := flag.Float64("size", 12, "Point size for font")
	flag.Parse()
	if *fontFile == "" {
		panic("Font file is required")
	}

	// Create a gRPC connection
	connection, err := grpc.Dial(fmt.Sprintf(":%d", clientPort), grpc.WithInsecure())
	errorHandler(err)
	// Create a flipdot client
	flipClient := flipdot.NewFlipdotClient(connection)

	// Create a flipdot controller
	flipdot, err := flipdot.NewFlipdot(flipClient)
	errorHandler(err)
	// Create a flipdot server
	server := flipapps.NewFlipappsServer(flipdot, readFont(*fontFile, *fontSize))
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	flipapps.RegisterFlipAppsServer(grpcServer, server)
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	errorHandler(err)
	// Start the server
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// Load font from disk
func readFont(filename string, size float64) font.Face {
	file, err := filepath.Abs(filename)
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	errorHandler(err)
	// Create the font face from the file
	face, err := text.NewFace(data, size)
	errorHandler(err)
	return face
}

// Generic error handler
func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
