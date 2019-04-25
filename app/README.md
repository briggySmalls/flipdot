# jim's-magic-sign ~ app

Business logic for the smart flipdot display

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

## Usage

The application binary has two subcommands: `flipapp app` and `flipapp mock`.