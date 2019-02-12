ENTRYPOINT=main.go
BIN_DIR=bin
BINARY=$(BIN_DIR)/flipcli

# Cross-compilation
PIOS=linux
PIARCH=arm
PIARM=7

# Generated protobufs
PROTO_DIR=flipdot
PROTO_FILE=$(PROTO_DIR)/flipdot.proto
PROTO_OUT=$(PROTO_DIR)/flipdot.pb.go

# Generated mocks
MOCKED_CLASS=FlipdotClient
MOCK_DIR=mock_flipdot
MOCK_FILE=$(MOCK_DIR)/flipdot.go

flipcli: protobuf
	go build -o $(BINARY) $(ENTRYPOINT)

flipcli-rpi: protobuf
	GOOS=$(PIOS) GOARCH=$(PIARCH) GOARM=$(PIARM) go build -a -o $(BINARY) $(ENTRYPOINT)

protobuf:
	protoc -I $(PROTO_DIR) $(PROTO_FILE) --go_out=plugins=grpc:$(PROTO_DIR)

test: mocks
	go test

mocks: protobuf
	mkdir -p $(MOCK_DIR)
	mockgen -source $(PROTO_OUT) -mock_names FlipdotClient=MockFlipdotClient > $(MOCK_FILE)

format:
	gofmt -w .

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_OUT)
