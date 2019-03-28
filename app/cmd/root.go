// Copyright © 2019 Sam Briggs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/briggySmalls/flipdot/app/flipapps"
	"github.com/briggySmalls/flipdot/app/flipdot"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio"
)

const (
	buttonDebounceDuration = time.Millisecond * 50
)

var cfgFile string

type config struct {
	clientAddress     string
	serverAddress     string
	fontFile          string
	fontSize          float64
	frameDurationSecs int
	appSecret         string
	appPassword       string
	buttonPin         uint8
	ledPin            uint8
	statusImage       string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flipapp",
	Short: "Application to display clock and messages on flipdot displays",
	Run: func(cmd *cobra.Command, args []string) {
		// Pull out config (from args/env/config file)
		config := validateConfig()
		// Create a client
		client, err := createClient(config.clientAddress)
		errorHandler(err)
		// Create a flipdot controller
		flipdot, err := flipdot.NewFlipdot(
			client,
			time.Duration(config.frameDurationSecs)*time.Second)
		errorHandler(err)
		// Create a button manager
		err = rpio.Open()
		errorHandler(err)
		defer rpio.Close()
		bm := createButtonManager(config.ledPin, config.buttonPin)
		// Get font
		font, err := readFont(config.fontFile, config.fontSize)
		// Create imager
		width, height := flipdot.Size()
		imager, err := createImager(config.statusImage, font, width, height, uint(len(flipdot.Signs())))
		errorHandler(err)
		// Create and start application
		app := flipapps.NewApplication(flipdot, bm, imager)
		go app.Run(time.Minute)
		// Create a flipapps server
		server := createServer(config.appSecret, config.appPassword, app.GetMessagesChannel(), flipdot.Signs())
		// Run server
		// Create a listener on TCP port
		lis, err := net.Listen("tcp", fmt.Sprintf(config.serverAddress))
		errorHandler(err)
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Init config
	cobra.OnInitialize(initConfig)

	// Always accept a config argument
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.flipapp.yaml)")

	// Define some root-command flags
	flags := rootCmd.Flags()
	flags.StringP("client-address", "c", "localhost:5001", "address used to connect to flipdot service")
	flags.StringP("server-address", "s", "0.0.0.0:5002", "address used to expose flipapp API over")
	flags.StringP("font-file", "f", "", "path to font .ttf file to display text with")
	flags.Float32P("font-size", "p", 0, "point size to obtain font face from font file")
	flags.Float32P("frame-duration", "d", 5, "Duration (in seconds) to display each frame of a message")
	flags.String("app-secret", "", "secret used to sign JWTs with")
	flags.String("app-password", "", "password required for authorisation")
	flags.Uint8("button-pin", 0, "GPIO pin that reads button state")
	flags.Uint8("led-pin", 0, "GPIO pin that illuminates button")
	flags.String("status-image", "", "Image to indicate new message status")

	// Add all flags to config
	viper.BindPFlags(flags)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".flipapp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".flipapp")
	}

	viper.AutomaticEnv() // read in environment variables that match
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer) // Separate environment variables with underscores

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Validate the supplied config
func validateConfig() config {
	clientAddress := viper.GetString("client-address")
	serverAddress := viper.GetString("server-address")
	fontFile := viper.GetString("font-file")
	fontSize := viper.GetFloat64("font-size")
	frameDuration := viper.GetInt("frame-duration")
	appSecret := viper.GetString("app-secret")
	appPassword := viper.GetString("app-password")
	buttonPin := viper.GetInt("button-pin")
	ledPin := viper.GetInt("led-pin")
	statusImage := viper.GetString("status-image")

	if serverAddress == "" {
		errorHandler(fmt.Errorf("server-address cannot be: %s", serverAddress))
	}
	if clientAddress == "" {
		errorHandler(fmt.Errorf("client-address cannot be: %s", clientAddress))
	}
	if fontSize == 0 {
		errorHandler(fmt.Errorf("font-size cannot be: %f", fontSize))
	}
	if fontFile == "" {
		errorHandler(fmt.Errorf("font-file cannot be: %s", fontFile))
	}
	if appSecret == "" {
		errorHandler(fmt.Errorf("app-secret cannot be: %s", appSecret))
	}
	if appPassword == "" {
		errorHandler(fmt.Errorf("app-password cannot be: %s", appPassword))
	}
	if statusImage == "" {
		errorHandler(fmt.Errorf("status-iamge cannot be: %s", statusImage))
	}

	fmt.Println("")
	fmt.Println("Starting server with the following configuration:")
	fmt.Printf("client-address: %s\n", clientAddress)
	fmt.Printf("server-address: %s\n", serverAddress)
	fmt.Printf("font-file: %s\n", fontFile)
	fmt.Printf("font-size: %f\n", fontSize)
	fmt.Printf("frame-duration: %d\n", frameDuration)
	fmt.Printf("button-pin: %d\n", buttonPin)
	fmt.Printf("led-pin: %d\n", ledPin)
	fmt.Printf("status-image: %s\n", statusImage)

	return config{
		clientAddress:     clientAddress,
		serverAddress:     serverAddress,
		fontFile:          fontFile,
		fontSize:          fontSize,
		frameDurationSecs: frameDuration,
		appSecret:         appSecret,
		appPassword:       appPassword,
		buttonPin:         uint8(buttonPin),
		ledPin:            uint8(ledPin),
		statusImage:       statusImage,
	}
}

// Generic error handler
func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
