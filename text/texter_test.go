package text

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/image/font"
)

func TestNewFace(t *testing.T) {
	// Load a font from disk
	file, err := filepath.Abs("m3x6.ttf")
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	// Test the NewFace function
	_, err = NewFace(data, 5)
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
		{"this, is. a: string", []string{"this", "is", "a", "string"}},
	}

	for _, table := range tables {
		calc := splitWords(table.input)
		if !reflect.DeepEqual(calc, table.output) {
			t.Errorf("SplitWords failed to split %s", table.input)
		}
	}
}

func TestToToLines(t *testing.T) {
	// Get test drawer
	f := getTestFont()
	d, err := createDrawer(f)
	errorHandler(err)
	tb := textBuilder{84, 7, f}

	// Prepare test table
	tables := []struct {
		input  string
		output []string
	}{
		{"Single line", []string{"Single line"}},
		{"This is a multi-line string", []string{"This is a multi-line", "string"}},
		{
			"This is a really really long string, maybe; it's even four lines",
			[]string{"This is a really", "really long string", "maybe it's even four", "lines"},
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
	// Get a font
	f := getTestFont()
	// Create the text builder
	tb := NewTextBuilder(84, 7, f)
	images, err := tb.Images("Hello my name is Sam. How's tricks?")
	if err != nil {
		t.Errorf("Image conversion returned error %s", err)
	}
	if len(images) != 1 {
		t.Errorf("Incorrect number of images: %d", len(images))
	}
	for _, img := range images {
		printImage(img)
	}
}

func getTestFont() font.Face {
	// Load a font from disk
	file, err := filepath.Abs("Smirnof.ttf")
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	// Test the NewFace function
	face, err := NewFace(data, 8)
	errorHandler(err)
	return face
}

func printImage(im image.Image) {
	c := color.Gray{0}
	for y := 0; y < im.Bounds().Dy(); y++ {
		fmt.Print("|")
		for x := 0; x < im.Bounds().Dx(); x++ {
			if im.At(x, y) == c {
				fmt.Print(" ")
			} else {
				fmt.Print("O")
			}
		}
		fmt.Println("|")
	}
}
