/*
This package implements an application microservice that sends instructions to a flipdot service
*/
package main

import (
	"fmt"

	"github.com/briggySmalls/flipcli/cmd"
	"github.com/briggySmalls/flipcli/flipdot"
	"google.golang.org/grpc"
)

func main() {
	cmd.Execute(createClient)
}

func createClient(port uint) (flipdot.FlipdotClient, *grpc.ClientConn, error) {
	// Get a gRPC connection
	connection, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	// Create the gRPC client for subcommands to use
	return flipdot.NewFlipdotClient(connection), connection, nil
}
