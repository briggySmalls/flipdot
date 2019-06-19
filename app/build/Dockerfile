## BUILD
FROM golang:alpine AS builder
# Install git for fetching the dependencies
RUN apk update && apk add --no-cache git make protobuf

WORKDIR /src/flipdot/app

# Fetch dependencies
COPY ./app/go.mod /src/flipdot/app/go.mod
COPY ./app/go.sum /src/flipdot/app/go.sum
RUN go mod download

# Copy the rest of the source
COPY . /src/flipdot

# Install the protobuf source generator binary
RUN go install github.com/golang/protobuf/protoc-gen-go

# Build the binary
RUN make build IS_PI=TRUE

# INSTALL
FROM scratch

# Copy in executable
COPY --from=builder /src/flipdot/app/bin/flipapp /go/bin/flipapp
# Copy in defaults for configuration
COPY ./app/config/config.yaml /app/
COPY ./app/assets/Smirnof.ttf /app/font.ttf
COPY ./app/assets/status.png /app/status.png

# Run the go program
ENTRYPOINT ["/go/bin/flipapp"]
CMD ["app", "--config", "/app/config.yaml"]
