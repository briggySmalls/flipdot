package cmd

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/briggySmalls/flipdot/app/flipapps"
	"github.com/briggySmalls/flipdot/app/flipdot"
	"github.com/briggySmalls/flipdot/app/text"
	"golang.org/x/image/font"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func createServer(appSecret, appPassword string, messagesIn chan flipapps.MessageRequest, signsInfo []*flipdot.GetInfoResponse_SignInfo) (grpcServer *grpc.Server) {
	grpcServer = flipapps.NewRpcServer(
		appSecret,
		appPassword,
		messagesIn,
		signsInfo,
	)
	// Register reflection service on gRPC server (for debugging).
	reflection.Register(grpcServer)
	return
}

func createClient(address string, isMock bool) (flipClient flipdot.FlipdotClient, err error) {
	if isMock {
		// Create a mock flipdot client
		return createMockFlipdotClient(), nil
	}
	// Create a gRPC connection to the remote flipdot server
	var connection *grpc.ClientConn
	connection, err = grpc.Dial(fmt.Sprintf(address), grpc.WithInsecure())
	if err != nil {
		return
	}
	// Create a flipdot client
	flipClient = flipdot.NewFlipdotClient(connection)
	return
}

func createImager(imageFile string, font font.Face, width, height, signCount uint) (imager flipapps.Imager, err error) {
	// Read in status image
	var statusImage image.Image
	statusImage, err = readImage(imageFile)
	if err != nil {
		return
	}
	// Create a text builder
	textBuilder := text.NewTextBuilder(width, height, font)
	// Create the imager
	imager = flipapps.NewImager(textBuilder, statusImage, signCount)
	return
}

// Load font from disk
func readFont(filename string, size float64) (face font.Face, err error) {
	var filePath string
	filePath, err = filepath.Abs(filename)
	if err != nil {
		return
	}
	var data []byte
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	// Create the font face from the file
	face, err = text.NewFace(data, size)
	if err != nil {
		return
	}
	return
}

func readImage(filename string) (image image.Image, err error) {
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
