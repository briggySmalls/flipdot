package text

import (
	"image"
	"image/color"
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
	images, err := tb.Images("Hello my name is Sam. How's tricks?", false)
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

// Helper function to get a font face
func getFont() (font font.Face) {
	return inconsolata.Regular8x16
}

// Check the individual pixels to see if any are set
func checkImage(im image.Image) bool {
	rect := im.Bounds()
	for x := 0; x < rect.Dx(); x++ {
		for y := 0; y < rect.Dy(); y++ {
			if (im.At(x, y) != color.RGBA{0, 0, 0, 0}) {
				return true
			}
		}
	}
	return false
}
