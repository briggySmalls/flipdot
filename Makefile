ENTRYPOINT=main.go
BIN_DIR=bin
BINARY=$(BIN_DIR)/flipcli

# Cross-compilation
PIOS=linux
PIARCH=arm
PIARM=7

# Generated protobufs
PROTO_DIR=protos
PROTO_SRCS=flipdot/flipdot.pb.go flipserver/flipserver.pb.go flipdot/flipdot.go
PROTO_BUFS=$(subst .go,.proto,$(PROTO_SRCS))
PROTO_MOCKS=$(subst .go,.mock.go,$(PROTO_SRCS))

# Generated mocks
MOCKED_CLASS=FlipdotClient
MOCK_DIR=mock_flipdot
MOCK_FILE=$(MOCK_DIR)/flipdot.go

flipcli: protobuf
	go build -o $(BINARY) $(ENTRYPOINT)

flipcli-rpi: protobuf
	GOOS=$(PIOS) GOARCH=$(PIARCH) GOARM=$(PIARM) go build -a -o $(BINARY) $(ENTRYPOINT)

protobuf: $(PROTO_SRCS)

%.pb.go: %.proto
	@echo Generating: $<
	protoc $(addprefix -I ,$(dir $(PROTO_BUFS))) $< --go_out=plugins=grpc:$(dir $<)

test: mocks
	go test

mocks: $(PROTO_MOCKS)

%.mock.go: %.go
	mockgen -source $< -package $(lastword $(subst /, ,$(dir $<))) > $@

format:
	gofmt -w .

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_SRCS) $(PROTO_MOCKS)
