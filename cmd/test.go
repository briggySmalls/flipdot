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
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/briggySmalls/flipcli/flipdot"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the signs",
	Long:  `Start or stop the test sign instruction`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start test signs",
	Long:  `Sends the start test signs instruction to all connected signs`,
	Run: func(cmd *cobra.Command, args []string) {
		if flipClient == nil {
			return
		}
		// Create timed context for request
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Send request
		response, err := flipClient.StartTest(ctx, &flipdot.StartTestRequest{})
		errorHandler(err)
		if response != nil {
			flipdotErrorHandler(*response.Error)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop test signs",
	Long:  `Sends the stop test signs instruction to all connected signs`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create timed context for request
		if flipClient == nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Send request
		response, err := flipClient.StopTest(ctx, &empty.Empty{})
		errorHandler(err)
		if response != nil {
			flipdotErrorHandler(*response.Error)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(startCmd)
	testCmd.AddCommand(stopCmd)
}
