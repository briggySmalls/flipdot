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

	"golang.org/x/image/font"

	"github.com/briggySmalls/flipcli/text"
	"github.com/spf13/cobra"
)

var fnt string

// textCmd represents the draw command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Write text to the signs",
	Args:  cobra.ExactArgs(1),
	Long: `Write text to the signs
The phrase is automatically wrapped, and the images staggered if necessary`,
	Run: func(cmd *cobra.Command, args []string) {
		// Send text
		err := controller.Text(args[0], getFont(fnt))
		errorHandler(err)
	},
}

func init() {
	rootCmd.AddCommand(textCmd)

	textCmd.Flags().StringVarP(&fnt, "font", "f", "", "Font to use to display text")
	textCmd.MarkFlagRequired("font")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// textCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// textCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getFont(fnt string) font.Face {
	// Load font from disk
	file, err := filepath.Abs(fnt)
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	errorHandler(err)
	// Create the font face from the file
	face, err := text.NewFace(data, 8)
	errorHandler(err)
	// Create a text builder
	return face
}
