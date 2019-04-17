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

package flipapp

import (
	"time"

	"github.com/briggySmalls/flipdot/app/internal/button"
	"github.com/spf13/cobra"
)

// mockCmd represents the mock command
var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "A mock version of the flipdot application",
	Long: `A terminal UI representation of the flipdot signs and the signals that get sent to it

Useful for debugging and development, especially when working with the web app.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config
		config := getMockConfig()
		// Create a mock flipdot client
		ui := createMockUI()
		// Assign client from UI
		client := ui
		// Create a button manager from UI
		bm := button.NewButtonManager(&ui.buttonPin, &ui.ledPin, time.Second, buttonDebounceDuration)

		// Run the rest of the app
		runApp(client, bm, config)
	},
}

func init() {
	rootCmd.AddCommand(mockCmd)
}

func getMockConfig() config {
	// Get the common config
	config := getCommonConfig()

	return config
}
