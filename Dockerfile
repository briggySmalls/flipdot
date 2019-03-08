## BUILD
FROM golang:alpine AS builder
# Install git for fetching the dependencies
RUN apk update && apk add --no-cache git make
WORKDIR $GOPATH/src/github.com/briggySmalls/flipapp
COPY . .
# Fetch dependencies using go get
RUN go get -d -v
# Build the binary
RUN make install IS_PI=TRUE

# INSTALL
FROM scratch
# Make some sensible default arguments
ENV CLIENT_PORT 5001
ENV SERVER_PORT 5002

# Copy in executable and font
COPY --from=builder /go/bin/flipapp /go/bin/flipapp
COPY ./Smirnoff.ttf /app

# Run the go program
ENTRYPOINT ["/go/bin/flipapps"]
CMD [ \
    "-client-port", "$CLIENT_PORT", \
    "-server-port", "$SERVER_PORT", \
    "-font", "/app/Smirnoff.ttf", \
    "-size", "6" \
]
