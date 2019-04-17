package text

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/shiny/text"
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

type TextBuilder interface {
	Images(text string, centre bool) ([]draw.Image, error)
}

func NewTextBuilder(width uint, height uint, font font.Face) TextBuilder {
	// Create and return a textBuilder
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

func (tb *textBuilder) Images(text string, centre bool) ([]draw.Image, error) {
	// Split the string up into lines
	lines, err := tb.toLines(text)
	errorHandler(err)

	// Create a drawer from the font
	d, err := createDrawer(tb.font)
	errorHandler(err)
	// Check font metrics
	m := d.Face.Metrics()
	charHeight := m.Ascent + m.Descent
	if charHeight.Floor() > int(tb.height) {
		return nil, fmt.Errorf("Font height %d larger than height %d", charHeight.Round(), tb.height)
	}
	// Draw the string
	var images []draw.Image
	for _, line := range lines {
		var xPos fixed.Int26_6 = 0
		if centre {
			lineWidth := d.MeasureString(line)
			xPos = (fixed.I(int(tb.width)) - lineWidth) / 2
		}
		// Reset the x position
		d.Dot = fixed.Point26_6{X: xPos, Y: m.Ascent}
		// Create a fresh destination
		d.Dst = image.NewGray(image.Rect(0, 0, int(tb.width), int(tb.height)))
		// Draw a new image
		d.DrawString(line)
		// Save the image
		images = append(images, d.Dst)
	}
	return images, nil
}

// Wrap text to multiple lines based off font and pixel width
func (tb *textBuilder) toLines(s string) ([]string, error) {
	// Create a frame
	var frame text.Frame
	frame.SetFace(tb.font)
	frame.SetMaxWidth(fixed.I(int(tb.width)))
	// Update the frame with the new text
	c := frame.NewCaret()
	c.WriteString(s)
	c.Close()
	f := &frame
	// Get the lines
	var lines []string
	for p := f.FirstParagraph(); p != nil; p = p.Next(f) {
		for l := p.FirstLine(f); l != nil; l = l.Next(f) {
			for b := l.FirstBox(f); b != nil; b = b.Next(f) {
				lines = append(lines, string(b.TrimmedText(f)[:]))
			}
		}
	}
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

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
