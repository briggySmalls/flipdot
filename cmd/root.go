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
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/briggySmalls/flipapp/flipapps"
	"github.com/briggySmalls/flipapp/flipdot"
	"github.com/briggySmalls/flipapp/text"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/image/font"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var cfgFile string

type config struct {
	clientAddress string
	serverAddress string
	fontFile      string
	fontSize      float64
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flipapp",
	Short: "Application to display clock and messages on flipdot displays",
	Run: func(cmd *cobra.Command, args []string) {
		// Root function
		config := validateConfig()

		// Create a gRPC connection
		connection, err := grpc.Dial(fmt.Sprintf(config.clientAddress), grpc.WithInsecure())
		errorHandler(err)
		// Create a flipdot client
		flipClient := flipdot.NewFlipdotClient(connection)

		// Create a flipdot controller
		flipdot, err := flipdot.NewFlipdot(flipClient)
		errorHandler(err)
		// Create a flipdot server
		server := flipapps.NewFlipappsServer(
			flipdot,
			readFont(config.fontFile, config.fontSize))
		// create a gRPC server object
		grpcServer := grpc.NewServer()
		// attach the Ping service to the server
		flipapps.RegisterFlipAppsServer(grpcServer, server)
		// Create a listener on TCP port
		lis, err := net.Listen("tcp", fmt.Sprintf(config.serverAddress))
		errorHandler(err)
		// Start the server
		// Register reflection service on gRPC server.
		reflection.Register(grpcServer)
		if err := grpcServer.Serve(lis); err != nil {
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

	// Add all flags to config
	viper.BindPFlag("client-address", flags.Lookup("client-address"))
	viper.BindPFlag("server-address", flags.Lookup("server-address"))
	viper.BindPFlag("font-file", flags.Lookup("font-file"))
	viper.BindPFlag("font-size", flags.Lookup("font-size"))
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

	fmt.Println("")
	fmt.Println("Starting server with the following configuration:")
	fmt.Printf("client-address: %s\n", clientAddress)
	fmt.Printf("server-address: %s\n", serverAddress)
	fmt.Printf("font-file: %s\n", fontFile)
	fmt.Printf("font-size: %f\n", fontSize)

	return config{
		clientAddress: clientAddress,
		serverAddress: serverAddress,
		fontFile:      fontFile,
		fontSize:      fontSize,
	}
}

// Load font from disk
func readFont(filename string, size float64) font.Face {
	file, err := filepath.Abs(filename)
	errorHandler(err)
	data, err := ioutil.ReadFile(file)
	errorHandler(err)
	// Create the font face from the file
	face, err := text.NewFace(data, size)
	errorHandler(err)
	return face
}

// Generic error handler
func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
