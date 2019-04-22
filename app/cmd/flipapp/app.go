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
	"fmt"
	"time"

	"github.com/briggySmalls/flipdot/app/internal/button"
	"github.com/briggySmalls/flipdot/app/internal/protos"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpio "github.com/stianeikeland/go-rpio/v4"
	"google.golang.org/grpc"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Flipdot application for deployment on Raspberry Pi",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get config
		config := getAppConfig()
		// Create a gRPC connection to the remote flipdot server
		connection, err := grpc.Dial(fmt.Sprintf(config.clientAddress), grpc.WithInsecure())
		errorHandler(err)
		// Create a flipdot client
		client := protos.NewDriverClient(connection)
		// Activate RPi GPIO
		err = rpio.Open()
		errorHandler(err)
		defer rpio.Close()
		// Create pins that interface with RPi GPIO
		ledPin := button.NewOutputPin(config.ledPin)
		buttonPin := button.NewTriggerPin(config.buttonPin)
		// Create a button manager
		bm := button.NewButtonManager(buttonPin, ledPin, time.Second, buttonDebounceDuration)

		// Run the rest of the app
		runApp(client, bm, config)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)

	flags := appCmd.Flags()
	flags.StringP("client-address", "c", "localhost:5001", "address used to connect to flipdot service")
	flags.Uint8("button-pin", 0, "GPIO pin that reads button state")
	flags.Uint8("led-pin", 0, "GPIO pin that illuminates button")
}

func getAppConfig() config {
	// First get the common config
	config := getCommonConfig()

	// Pull out more
	clientAddress := viper.GetString("client-address")
	buttonPin := viper.GetInt("button-pin")
	ledPin := viper.GetInt("led-pin")

	// Validate additional config
	if clientAddress == "" {
		errorHandler(fmt.Errorf("client-address cannot be: %s", clientAddress))
	}

	// Print additional app config
	fmt.Printf("APP CONFIG")
	fmt.Printf("client-address: %s\n", clientAddress)
	fmt.Printf("button-pin: %d\n", buttonPin)
	fmt.Printf("led-pin: %d\n", ledPin)

	// Update config
	config.clientAddress = clientAddress
	config.buttonPin = uint8(buttonPin)
	config.ledPin = uint8(ledPin)

	return config
}
