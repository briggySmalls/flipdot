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

func TestToToLines(t *testing.T) {
	// Get test font
	f := getFont()
	tb, ok := NewTextBuilder(140, 17, f).(*textBuilder)
	if !ok {
		t.Fatal("TextBuilder is not a textBuilder")
	}

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
		{"This string\nhas\nnewlines.", []string{"This string", "has", "newlines."}},
	}

	for _, table := range tables {
		lines, err := tb.toLines(table.input)
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
	tb := NewTextBuilder(120, 17, f)
	images, err := tb.Images("Hello my name is Sam. How's tricks?", false)
	if err != nil {
		t.Fatalf("Image conversion returned error: %s", err)
		return
	} else if len(images) != 3 {
		t.Fatalf("Incorrect number of images: %d", len(images))
	}
	for _, img := range images {
		if !checkImage(img) {
			t.Error("Image empty")
		}
	}
}

func TestCentring(t *testing.T) {
	// Get test font
	f := getFont()
	// Create the text builder
	var width uint = 20
	tb := NewTextBuilder(width, 17, f)
	// Write a vertical pipe (should be first pixels)
	images, err := tb.Images("|", false)
	if err != nil {
		t.Fatalf("Image conversion returned error: %s", err)
		return
	} else if len(images) != 1 {
		t.Fatalf("Incorrect number of images: %d", len(images))
	}
	// Check first pixels are populated
	if (images[0].At(0, f.Metrics().Ascent.Round()) == color.RGBA{0, 0, 0, 0}) {
		t.Error("Pixels not set in left-aligned image")
	}
	// Write a vertical pipe (should be middle pixels)
	images, err = tb.Images("|", true)
	if err != nil {
		t.Fatalf("Image conversion returned error: %s", err)
		return
	} else if len(images) != 1 {
		t.Fatalf("Incorrect number of images: %d", len(images))
	}
	// Check middle pixels are populated
	if (images[0].At(int(width)/2, f.Metrics().Ascent.Round()) == color.RGBA{0, 0, 0, 0}) {
		t.Error("Pixels not set in left-aligned image")
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
