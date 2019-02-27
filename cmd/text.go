// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/briggySmalls/flipcli/flipdot"

	"github.com/briggySmalls/flipcli/text"
	"github.com/spf13/cobra"
)

var font string

// textCmd represents the draw command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Write text to the signs",
	Args:  cobra.ExactArgs(1),
	Long: `Write text to the signs
The phrase is automatically wrapped, and the images staggered if necessary`,
	Run: func(cmd *cobra.Command, args []string) {
		signs := getSigns()
		err := checkSigns(signs)
		errorHandler(err)
		// Get a text builder
		tb := getTextBuilder(font)
		// Make the text images
		images, err := tb.Images(args[0])
		errorHandler(err)
		// Write the images
		var signNames []string
		for _, sign := range signs {
			signNames = append(signNames, sign.Name)
		}
		sendText(images, signNames)
	},
}

func init() {
	rootCmd.AddCommand(textCmd)

	textCmd.Flags().StringVarP(&font, "font", "f", "", "Font to use to display text")
	textCmd.MarkFlagRequired("font")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// textCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// textCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTextBuilder(font string) text.TextBuilder {
	// Load font from disk
	file, err := filepath.Abs(font)
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	errorHandler(err)
	// Create the font face from the file
	face, err := text.NewFace(data, 8)
	errorHandler(err)
	// Create a text builder
	return text.NewTextBuilder(84, 7, face)
}

func sendText(images []text.Image, signNames []string) {
	// Send any relevant images
	images = sendFrame(images, signNames)
	// Check if we need to go on
	if len(images) == 0 {
		return
	}
	// Create a ticker
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	// Write images periodically
	for len(images) > 0 {
		select {
		case <-ticker.C:
			images = sendFrame(images, signNames)
		}
	}
}

// Send a set of images to available signs
func sendFrame(images []text.Image, signNames []string) (leftover []text.Image) {
	for _, sign := range signNames {
		// Stop sending if there are no more images left
		if len(images) == 0 {
			return images
		}
		// Pop an image off the stack and send it
		var image text.Image
		image, images = images[0], images[1:]
		writeImage(image, sign)
	}
	return images
}

// Write an image to the specified sign
func writeImage(image text.Image, sign string) {
	// Send request
	ctx, cancel := getContext()
	defer cancel()
	_, err := flipClient.Draw(ctx, &flipdot.DrawRequest{
		Sign:  sign,
		Image: image.Slice(),
	})
	errorHandler(err)
}

// Check that all signs have the same width/height
func checkSigns(signs []*flipdot.GetInfoResponse_SignInfo) error {
	var width, height int32
	for i, sign := range signs {
		if i == 0 {
			width = sign.Width
			height = sign.Height
		} else {
			if width != sign.Width {
				return fmt.Errorf("Sign width %d != %d", sign.Width, width)
			} else if height != sign.Height {
				return fmt.Errorf("Sign height %d != %d", sign.Height, height)
			}
		}
	}
	return nil
}

// Get the list of signs from the client
func getSigns() (signs []*flipdot.GetInfoResponse_SignInfo) {
	// Get the signs
	context, cancel := getContext()
	defer cancel()
	response, err := flipClient.GetInfo(context, &flipdot.GetInfoRequest{})
	errorHandler(err)
	return response.Signs
}
