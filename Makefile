ENTRYPOINT=./main.go
BIN_DIR=bin
BINARY=$(BIN_DIR)/flipapp

# Cross-compilation
PIOS=linux
PIARCH=arm
PIARM=7

# Generated protobufs
PROTO_DIR=protos
PROTO_SRCS=flipdot/flipdot.pb.go flipapps/flipapps.pb.go
PROTO_BUFS=$(subst .pg.go,.proto,$(PROTO_SRCS))
MOCKS=$(subst .go,.mock.go,$(PROTO_SRCS)) flipdot/flipdot.mock.go

# Generated mocks
MOCKED_CLASS=FlipdotClient
MOCK_DIR=mock_flipdot
MOCK_FILE=$(MOCK_DIR)/flipdot.go

flipapps: protobuf
	go build -o $(BINARY) $(ENTRYPOINT)

flipapps-rpi: protobuf
	GOOS=$(PIOS) GOARCH=$(PIARCH) GOARM=$(PIARM) go build -a -o $(BINARY) $(ENTRYPOINT)

protobuf: $(PROTO_SRCS)

%.pb.go: %.proto
	@echo Generating: $<
	protoc $(addprefix -I ,$(dir $(PROTO_BUFS))) --go_out=plugins=grpc:../../.. $<

test: mocks
	go test ./...

mocks: $(MOCKS)

%.mock.go: %.go
	mockgen -source $< -package $(lastword $(subst /, ,$(dir $<))) > $@

format:
	gofmt -w .

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_SRCS) $(MOCKS)
