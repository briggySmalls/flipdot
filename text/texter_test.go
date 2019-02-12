package main

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

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

func TestToOneLine(t *testing.T) {
	// Get test drawer
	f := getTestFont()
	d, err := createDrawer(f)
	errorHandler(err)
	// Test toLines
	tb := textBuilder{84, 7}
	lines, err := tb.toLines(*d, "Single line")
	errorHandler(err)
	// Assert
	if len(lines) != 1 {
		t.Errorf("Unexpected number of lines: %d", len(lines))
	}
}

func TestImages(t *testing.T) {
	// Get a font
	f := getTestFont()
	// Create the text builder
	tb := NewTextBuilder(84, 7)
	images, err := tb.Images("Turn into image", f)
	if err != nil {
		t.Errorf("Image conversion returned error %s", err)
	}
	if len(images) != 1 {
		t.Errorf("Incorrect number of images: %d", len(images))
	}
}

func getTestFont() []byte {
	// Get a font
	file, err := filepath.Abs("fonts/Nintendo-Entertainment-System/Nintendo Entertainment System.ttf")
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	errorHandler(err)
	return data
}
