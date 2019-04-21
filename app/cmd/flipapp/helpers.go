package flipapp

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/imaging"
	"github.com/briggySmalls/flipdot/app/internal/protos"
	"github.com/briggySmalls/flipdot/app/internal/server"
	"github.com/briggySmalls/flipdot/app/internal/text"
	"golang.org/x/image/font"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func createServer(appSecret, appPassword string, tokenExpiry time.Duration, messagesIn chan protos.MessageRequest, signsInfo []*protos.GetInfoResponse_SignInfo) (grpcServer *grpc.Server) {
	grpcServer = server.NewRpcServer(
		appSecret,
		appPassword,
		tokenExpiry,
		messagesIn,
		signsInfo,
	)
	// Register reflection service on gRPC server (for debugging).
	reflection.Register(grpcServer)
	return
}

func createImager(imageFile string, font font.Face, width, height, signCount uint) (imager imaging.Imager, err error) {
	// Read in status image
	var statusImage image.Image
	statusImage, err = readImage(imageFile)
	if err != nil {
		return
	}
	// Create a text builder
	textBuilder := text.NewTextBuilder(width, height, font)
	// Create the imager
	imager = imaging.NewImager(textBuilder, statusImage, signCount)
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
