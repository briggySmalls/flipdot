package text

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func NewFace(data []byte, points float64) (font.Face, error) {
	// Parse the font
	font, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	// Turn into a face
	opts := truetype.Options{
		Size: points,
	}
	return truetype.NewFace(font, &opts), nil
}

func Slice(image image.Image) []bool {
	bgColor := color.Gray{0}
	// Create an array for the image
	rows := image.Bounds().Dy()
	cols := image.Bounds().Dx()
	binImage := make([]bool, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			binImage[r*cols+c] = image.At(c, r) != bgColor
		}
	}
	return binImage
}

func UnSlice(data []bool, width, height int) (img image.Image, err error) {
	if width*height != len(data) {
		err = fmt.Errorf("Width %d, height %d incompatible with data length %d",
			width, height, len(data))
		return
	}
	grey := image.NewGray(image.Rect(0, 0, width, height))
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			var colour color.Gray
			if data[r*width+c] {
				colour = color.Gray{255}
			} else {
				colour = color.Gray{0}
			}
			grey.SetGray(c, r, colour)
		}
	}
	return
}

type TextBuilder interface {
	Images(text string) ([]image.Image, error)
}

func NewTextBuilder(width uint, height uint, font font.Face) TextBuilder {
	return &textBuilder{
		width:  width,
		height: height,
		font:   font,
	}
}

type textBuilder struct {
	width  uint
	height uint
	font   font.Face
}

func (tb *textBuilder) Images(text string) ([]image.Image, error) {
	// Create a drawer from the font
	d, err := createDrawer(tb.font)
	errorHandler(err)
	// Check font metrics
	m := d.Face.Metrics()
	charHeight := m.Ascent + m.Descent
	if charHeight.Floor() > int(tb.height) {
		return nil, fmt.Errorf("Font height %d larger than height %d", charHeight.Round(), tb.height)
	}
	// Split the string up into lines
	lines, err := tb.toLines(*d, text)
	errorHandler(err)
	// Draw the string
	var images []image.Image
	for _, line := range lines {
		// Reset the x position
		d.Dot = fixed.Point26_6{X: 0, Y: m.Ascent}
		// Create a fresh destination
		d.Dst = image.NewGray(image.Rect(0, 0, int(tb.width), int(tb.height)))
		// Draw a new image
		d.DrawString(line)
		// Save the image
		images = append(images, d.Dst)
	}
	return images, nil
}

func (tb *textBuilder) toLines(d font.Drawer, s string) ([]string, error) {
	words := splitWords(s)
	var lines []string
	start := 0
	for end, word := range words {
		// Build a query line
		queryLine := strings.Join(words[start:end+1], " ")
		// Determine if string fits
		if d.MeasureString(queryLine) > fixed.I(int(tb.width)) {
			// String doesn't fit on line
			if end == 0 {
				// Single word is too big for a line
				return nil, fmt.Errorf("Word '%s' too large to fit on line of width %d", word, tb.width)
			}
			// The previous must have fit
			lines = append(lines, strings.Join(words[start:end], " "))
			// Start again
			start = end
		}
	}
	// Add any remaining lines
	lines = append(lines, strings.TrimSpace(strings.Join(words[start:len(words)], " ")))
	return lines, nil
}

func createDrawer(face font.Face) (*font.Drawer, error) {
	// Establish the baseline
	m := face.Metrics()
	// Update the drawer with font-specific fields
	return &font.Drawer{
		Src:  &image.Uniform{color.Gray{255}},
		Face: face,
		Dot:  fixed.Point26_6{X: 0, Y: m.Ascent},
	}, nil
}

// Helper function to split words up
func splitWords(s string) []string {
	return strings.Fields(s)
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
