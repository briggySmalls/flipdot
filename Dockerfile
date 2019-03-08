## BUILD
FROM golang:alpine AS builder
# Install git for fetching the dependencies
RUN apk update && apk add --no-cache git make protobuf
WORKDIR $GOPATH/src/github.com/briggySmalls/flipapp

# Get the protobuf source generator (do this prior to copy)
RUN go get -u github.com/golang/protobuf/protoc-gen-go
# Fetch other dependencies using go get
COPY . .
RUN go get -d -v

# Build the binary
RUN make build IS_PI=TRUE

# INSTALL
FROM scratch
# Make some sensible default arguments
ENV CLIENT_PORT 5001
ENV SERVER_PORT 5002

# Copy in executable and font
COPY --from=builder /go/src/github.com/briggySmalls/flipapp/bin/flipapp /go/bin/flipapp
COPY ./Smirnof.ttf /app

# Run the go program
ENTRYPOINT ["/go/bin/flipapps"]
CMD [ \
    "--client-port", "$CLIENT_PORT", \
    "--server-port", "$SERVER_PORT", \
    "--font-file", "/app/Smirnof.ttf", \
    "--font-size", "6" \
]
