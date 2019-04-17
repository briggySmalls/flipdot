package imaging

import (
	"image"
	"image/color"
	"image/draw"
	reflect "reflect"
	"testing"
	"time"

	"github.com/briggySmalls/flipdot/app/text"
	gomock "github.com/golang/mock/gomock"
)

// Test converting a 2D image into a slice
func TestSlice(t *testing.T) {
	// Prepare test table
	tables := []struct {
		input  image.Image
		output []bool
	}{
		{createTestImage(color.Gray{1}, image.Rect(0, 0, 2, 3)), []bool{true, true, true, true, true, true}},
		{createTestImage(color.Gray{0}, image.Rect(0, 0, 3, 3)), []bool{false, false, false, false, false, false, false, false, false}},
	}
	for _, table := range tables {
		slice := Slice(table.input)
		if !reflect.DeepEqual(slice, table.output) {
			t.Errorf("Image %s not sliced correctly", table.input)
		}
	}
}

// Test converting the time into images
func TestClock(t *testing.T) {
	// Create the test objects
	imgr, tb := createImagerTestObjects(t, 1, 1, nil)

	// Expect a call to create images from text
	tb.EXPECT().Images("Tue 1 Jan\n12:30 pm", true).Return([]draw.Image{
		image.NewGray(image.Rect(0, 0, 1, 1)),
	}, nil)
	// Request time be drawn, without status image
	images, err := imgr.Clock(time.Date(2019, 1, 1, 12, 30, 0, 0, time.UTC), false)
	if err != nil {
		t.Fatal(err)
	}
	// Check images are expected and have no added content
	if len(images) != 1 {
		t.Fatalf("Unexpected number of images %d", len(images))
	}
	if images[0].Data[0] {
		t.Errorf("Image unexpectedly modified")
	}
}

// Test converting the time with a status
func TestClockStatus(t *testing.T) {
	// Create the test objects
	statusImage := createTestImage(color.Gray{255}, image.Rect(0, 0, 1, 1))
	imgr, tb := createImagerTestObjects(t, 2, 1, statusImage)

	// Expect a call to create images from text
	tb.EXPECT().Images("Tue 1 Jan\n12:30 pm", true).Return([]draw.Image{
		image.NewGray(image.Rect(0, 0, 2, 1)),
	}, nil)
	// Request time be drawn, without status image
	images, err := imgr.Clock(time.Date(2019, 1, 1, 12, 30, 0, 0, time.UTC), true)
	if err != nil {
		t.Fatal(err)
	}
	// Check images are expected and have no added content
	if len(images) != 1 {
		t.Fatalf("Unexpected number of images %d", len(images))
	}
	if images[0].Data[0] {
		t.Error("First pixel was set non-zero")
	}
	if !images[0].Data[1] {
		t.Error("Status image not set")
	}
}

func TestMessage(t *testing.T) {
	// Create test objects
	imgr, tb := createImagerTestObjects(t, 1, 1, nil)

	// Expect call to Images, and return dummy image
	fakeImages := []draw.Image{image.NewGray(image.Rect(0, 0, 1, 1))}
	gomock.InOrder(
		tb.EXPECT().Images("From: Sam", true).Return(fakeImages, nil),
		tb.EXPECT().Images("", false).Return(fakeImages, nil),
		tb.EXPECT().Images("hello", true).Return(fakeImages, nil),
	)
	// Call message
	images, err := imgr.Message("Sam", "hello")
	if err != nil {
		t.Fatal(err)
	}
	// Assert images are as expected
	if len(images) != 3 {
		t.Error("Complete frames not sent")
	}
}

func createImagerTestObjects(t *testing.T, width, height int, statusImage image.Image) (imager Imager, tb *text.MockTextBuilder) {
	// Create a mock textbuilder
	ctrl := gomock.NewController(t)
	tb = text.NewMockTextBuilder(ctrl)
	// Create an imager
	imager = NewImager(tb, statusImage, 2)
	return
}

func createTestImage(c color.Color, r image.Rectangle) image.Image {
	// Draw the image
	dst := image.NewGray(r)
	draw.Draw(dst, dst.Bounds(), image.NewUniform(c), image.Point{X: 0, Y: 0}, draw.Src)
	return dst
}
