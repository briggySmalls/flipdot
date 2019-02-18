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
		// Get a text builder
		tb := getTextBuilder(font)
		// Make the text images
		images, err := tb.Images(args[0])
		errorHandler(err)
		// Write the images
		sendText(images)
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

func sendText(images []text.Image) {
	// Create a ticker
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	// Write images periodically
	counter := 0
	for {
		select {
		case <-ticker.C:
			// Write to top sign
			if tryWriteImage(images, counter, "top") {
				counter++
			} else {
				return
			}
			// Write to bottom sign
			if tryWriteImage(images, counter, "bottom") {
				counter++
			} else {
				return
			}
		}
	}
}

func tryWriteImage(images []text.Image, index int, sign string) bool {
	if index >= len(images) {
		return false
	}
	// Send request
	ctx, cancel := getContext()
	defer cancel()
	flipClient.Draw(ctx, &flipdot.DrawRequest{
		Sign:  sign,
		Image: images[index].Slice(),
	})
	return true
}
