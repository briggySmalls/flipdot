## BUILD
FROM golang:alpine AS builder
# Install git for fetching the dependencies
RUN apk update && apk add --no-cache git make protobuf

# Get the protobuf source generator (do this prior to copy)
RUN go get -u github.com/golang/protobuf/protoc-gen-go
# Fetch other dependencies using go get
COPY . $GOPATH/src/github.com/briggySmalls/flipdot
WORKDIR $GOPATH/src/github.com/briggySmalls/flipdot/app
RUN go get -d -v

# Build the binary
RUN make build IS_PI=TRUE

# INSTALL
FROM scratch

# Copy in executable
COPY --from=builder /go/src/github.com/briggySmalls/flipdot/app/bin/flipapp /go/bin/flipapp
# Copy in defaults for configuration
COPY ./app/config.yaml /app/
COPY ./app/Smirnof.ttf /app/font.ttf
COPY ./app/bell.png /app/status.png

# Run the go program
ENTRYPOINT ["/go/bin/flipapp"]
CMD ["app", "--config", "/app/config.yaml"]
