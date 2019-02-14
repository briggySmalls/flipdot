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
		DPI:  72 / 64,
	}
	return truetype.NewFace(font, &opts), nil
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
	if charHeight > fixed.Int26_6(tb.height) {
		return nil, fmt.Errorf("Font height %d larger than height %d", charHeight, tb.height)
	}
	// Split the string up into lines
	lines, err := tb.toLines(*d, text)
	errorHandler(err)
	// Draw the string
	var images []image.Image
	for _, line := range lines {
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
	currentLine := strings.Builder{}
	var previous string
	for i, word := range words {
		// Keep string prior to addition
		previous = currentLine.String()
		// Add new word
		currentLine.WriteString(word)
		// Determine if string fits
		width := d.MeasureString(currentLine.String())
		if width > fixed.Int26_6(tb.width) {
			// String doesn't fit on line
			if i == 0 {
				// Single word is too big for a line
				return nil, fmt.Errorf("Word %s too large", word)
			}
			// Add previous (fitting) line to lines
			lines = append(lines, previous)
			// Start again
			currentLine.Reset()
		}
	}
	// Add any remaining lines
	lines = append(lines, previous)
	return lines, nil
}

func createDrawer(face font.Face) (*font.Drawer, error) {
	// Establish the baseline
	m := face.Metrics()
	// Update the drawer with font-specific fields
	return &font.Drawer{
		Src:  &image.Uniform{color.Gray{255}},
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(0), Y: m.Ascent},
	}, nil
}

func splitWords(s string) []string {
	s = strings.Map(removePunctuation, s)
	return strings.Fields(s)
}

func removePunctuation(r rune) rune {
	if strings.ContainsRune(".,:;", r) {
		return -1
	} else {
		return r
	}
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
