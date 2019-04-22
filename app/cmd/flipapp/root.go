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

package flipapp

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/briggySmalls/flipdot/app/internal"
	"github.com/briggySmalls/flipdot/app/internal/button"
	"github.com/briggySmalls/flipdot/app/internal/client"
	"github.com/briggySmalls/flipdot/app/internal/protos"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	tokenExpiry       time.Duration
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flipapp",
	Short: "Application to display clock and messages on flipdot displays",
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
	persistentFlags := rootCmd.PersistentFlags()
	persistentFlags.StringP("server-address", "s", "0.0.0.0:5002", "address used to expose flipapp API over")
	persistentFlags.StringP("font-file", "f", "", "path to font .ttf file to display text with")
	persistentFlags.Float32P("font-size", "p", 0, "point size to obtain font face from font file")
	persistentFlags.Float32P("frame-duration", "d", 5, "duration (in seconds) to display each frame of a message")
	persistentFlags.String("app-secret", "", "secret used to sign JWTs with")
	persistentFlags.String("app-password", "", "password required for authorisation")
	persistentFlags.String("status-image", "", "image to indicate new message status")
	persistentFlags.DurationP("token-expiry", "t", time.Hour, "duration after which a login token expires")

	// Add all flags to config
	viper.BindPFlags(persistentFlags)
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
func getCommonConfig() config {
	serverAddress := viper.GetString("server-address")
	fontFile := viper.GetString("font-file")
	fontSize := viper.GetFloat64("font-size")
	frameDuration := viper.GetInt("frame-duration")
	appSecret := viper.GetString("app-secret")
	appPassword := viper.GetString("app-password")
	statusImage := viper.GetString("status-image")
	tokenExpiry := viper.GetDuration("token-expiry")

	if serverAddress == "" {
		errorHandler(fmt.Errorf("server-address cannot be: %s", serverAddress))
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
		errorHandler(fmt.Errorf("status-image cannot be: %s", statusImage))
	}
	if tokenExpiry == 0 {
		errorHandler(fmt.Errorf("token-expiry cannot be: %d", tokenExpiry))
	}

	fmt.Println("")
	fmt.Println("Starting server with the following configuration:")
	fmt.Printf("server-address: %s\n", serverAddress)
	fmt.Printf("font-file: %s\n", fontFile)
	fmt.Printf("font-size: %f\n", fontSize)
	fmt.Printf("frame-duration: %d\n", frameDuration)
	fmt.Printf("status-image: %s\n", statusImage)
	fmt.Printf("token-expiry: %d\n", tokenExpiry)

	return config{
		serverAddress:     serverAddress,
		fontFile:          fontFile,
		fontSize:          fontSize,
		frameDurationSecs: frameDuration,
		appSecret:         appSecret,
		appPassword:       appPassword,
		statusImage:       statusImage,
		tokenExpiry:       tokenExpiry,
	}
}

// Create components and run application
func runApp(clnt protos.DriverClient, bm button.ButtonManager, config config) {
	// Create a flipdot controller
	flippy, err := client.NewFlipdot(
		clnt,
		time.Duration(config.frameDurationSecs)*time.Second)
	errorHandler(err)

	// Get font
	font, err := readFont(config.fontFile, config.fontSize)
	// Create imager
	width, height := flippy.Size()
	imager, err := createImager(config.statusImage, font, width, height, uint(len(flippy.Signs())))
	errorHandler(err)

	// Create and start application
	app := internal.NewApplication(flippy, bm, imager)
	go app.Run(30 * time.Second)
	// Create a flipapps server
	server := createServer(config.appSecret, config.appPassword, config.tokenExpiry, app.GetMessagesChannel(), flippy.Signs())
	// Run server
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(config.serverAddress))
	errorHandler(err)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// Generic error handler
func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
