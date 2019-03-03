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

func createClient(port uint) (flippy flipdot.Flipdot, connection *grpc.ClientConn, err error) {
	// Get a gRPC connection
	connection, err = grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		return
	}
	// Create the gRPC client for subcommands to use
	client := flipdot.NewFlipdotClient(connection)
	flippy, err = flipdot.NewFlipdot(client)
	return
}
