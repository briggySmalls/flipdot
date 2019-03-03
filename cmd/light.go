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

	"github.com/spf13/cobra"
)

// lightCmd represents the light command
var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("light called")
	},
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turns on the lights",
	Long:  `Turns on the lights that illuminate the flipdot displays`,
	Run: func(cmd *cobra.Command, args []string) {
		// Send light on instruction
		err := controller.LightOn()
		errorHandler(err)
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turns off the lights",
	Long:  `Turns off the lights that illuminate the flipdot displays`,
	Run: func(cmd *cobra.Command, args []string) {
		// Send light off instruction
		err := controller.LightOff()
		errorHandler(err)
	},
}

func init() {
	rootCmd.AddCommand(lightCmd)

	lightCmd.AddCommand(onCmd)
	lightCmd.AddCommand(offCmd)
}
