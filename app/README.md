# jim's-magic-sign ~ app

Business logic for the smart flipdot display, written in Go.

## Features

- Makes calls to [driver service](../driver) to display application
- Defaults to displaying time/date
- Exposes gRPC 'service' to enqueue messages
- Flashes button when messages are in the queue
- Listens for button press to display queued messages

## Installation

As a Go module, installation is as simple as:

```bash
go mod download
```

You need to perform this step, rather than just building or testing immediately, in order to download the gRPC protobuf source generator.

Now you can build or test the application:

```
# Build app
make

# Test app
make test
```

A docker container can be built using the following:
```
make docker
```

The built version of this container can be found on dockerhub as `briggysmalls/flipapp`.

## Usage

The application binary has two subcommands: `app` and `mock`.

The `app` subcommand is the main software, expecting to connect to a gRPC driver service and raspberry Pi hardware.

The `mock` subcommand is a version of the application that mocks away the gRPC driver service, instead simulating the signs on the console.
This command is useful for quick development of the build of the application, as well as providing a stubbed backend for the [web](../web) project.