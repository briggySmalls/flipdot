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

# Copy in executable and font
COPY --from=builder /go/src/github.com/briggySmalls/flipdot/app/bin/flipapp /go/bin/flipapp
COPY ./app/Smirnof.ttf /app/
COPY ./app/bell.png /app/
COPY ./app/config.yaml /app/

# Run the go program
ENTRYPOINT ["/go/bin/flipapp"]
CMD ["--config", "/app/config.yaml"]

