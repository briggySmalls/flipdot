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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var port *uint
var clientFactory func(uint) (flipdot.Flipdot, *grpc.ClientConn, error)
var controller flipdot.Flipdot
var connection *grpc.ClientConn

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flipcli",
	Short: "Simple CLI for testing the flipdriver service",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Create the client, with the specified connection
		var err error
		controller, connection, err = clientFactory(*port)
		errorHandler(err)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if connection != nil {
			// Close the connection if it was ever created
			connection.Close()
		}
	},
}

func init() {
	port = rootCmd.PersistentFlags().Uint("port", 5001, "Port of the gRPC server")
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func flipdotErrorHandler(err flipdot.Error) {
	if err.Code != 0 {
		panic(err.Message)
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func Execute(cf func(uint) (flipdot.Flipdot, *grpc.ClientConn, error)) {
	// Keep hold of the factory method
	clientFactory = cf
	// Execute adds all child commands to the root command and sets flags appropriately.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
