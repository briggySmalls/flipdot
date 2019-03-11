package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/inconsolata"
)

func TestNewFace(t *testing.T) {
	// Load a font from disk
	data := goregular.TTF
	// Test the NewFace function
	_, err := NewFace(data, 5)
	if err != nil {
		t.Error(err)
	}
}

func TestSplitWords(t *testing.T) {
	tables := []struct {
		input  string
		output []string
	}{
		{"this is a string", []string{"this", "is", "a", "string"}},
		{"this, is. a: string", []string{"this,", "is.", "a:", "string"}},
	}

	for _, table := range tables {
		calc := splitWords(table.input)
		if !reflect.DeepEqual(calc, table.output) {
			t.Errorf("SplitWords failed to split %s", table.input)
		}
	}
}

func TestToToLines(t *testing.T) {
	// Get test font
	f := getFont()
	// Get test drawer
	d, err := createDrawer(f)
	errorHandler(err)
	tb := textBuilder{140, 17, f}

	// Prepare test table
	tables := []struct {
		input  string
		output []string
	}{
		{"Single line", []string{"Single line"}},
		{"This is a multi-line string", []string{"This is a", "multi-line string"}},
		{
			"This is a really really long string, maybe; it's four lines",
			[]string{"This is a really", "really long", "string, maybe;", "it's four lines"},
		},
	}

	for _, table := range tables {
		lines, err := tb.toLines(*d, table.input)
		// Check operation passed
		if err != nil {
			t.Error(err)
		}
		// Check result is as expected
		if !reflect.DeepEqual(lines, table.output) {
			t.Errorf("toLines failed to split: %s", table.input)
		}
	}
}

func TestImages(t *testing.T) {
	// Get test font
	f := getFont()
	// Create the text builder
	var rowCount uint = 120
	tb := NewTextBuilder(rowCount, 17, f)
	images, err := tb.Images("Hello my name is Sam. How's tricks?")
	if err != nil {
		t.Fatalf("Image conversion returned error: %s", err)
		return
	}
	if len(images) != 3 {
		t.Fatalf("Incorrect number of images: %d", len(images))
	}
	for _, img := range images {
		if !checkImage(img) {
			t.Error("Image empty")
		}
	}
}

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
		bounds := table.input.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y
		unslice, err := UnSlice(slice, width, height)
		if err != nil {
			t.Error("Failed to unslice slice")
		}
		reflect.DeepEqual(unslice, table.input)
	}
}

// Helper function to get a font face
func getFont() (font font.Face) {
	return inconsolata.Regular8x16
}

func checkImage(im image.Image) bool {
	// Check the individual pixels to see if any are set
	for _, pix := range Slice(im) {
		if pix {
			return true
		}
	}
	return false
}

func createTestImage(c color.Color, r image.Rectangle) image.Image {
	// Draw the image
	dst := image.NewGray(r)
	draw.Draw(dst, dst.Bounds(), image.NewUniform(c), image.Point{X: 0, Y: 0}, draw.Src)
	return dst
}

// Helper function to print out an image on the command line
func printImage(images []image.Image, rowCount uint) {
	// Draw the image on the command line
	for _, img := range images {
		for i, pix := range Slice(img) {
			if pix {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
			if uint(i)%rowCount == rowCount-1 {
				fmt.Println("")
			}
		}
	}
}
