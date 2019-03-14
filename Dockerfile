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

# Copy in executable and font
COPY --from=builder /go/src/github.com/briggySmalls/flipapp/bin/flipapp /go/bin/flipapp
COPY ./Smirnof.ttf /app/
COPY ./config.yaml /app/

# Run the go program
ENTRYPOINT ["/go/bin/flipapp"]
CMD ["--config", "/app/config.yaml"]

