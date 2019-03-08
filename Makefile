# Build configuration
MAIN_FILE=./main.go
BIN_DIR=bin
BUILD_CMD?=go build
BUILD_ARGS?=-o $(BIN_DIR)/flipapp

# Cross-compilation
PI_OS=linux
PI_ARCH=arm
PI_ARM=7
IS_PI?=
ifdef IS_PI
	BUILD_ARGS+= -a
	ENVS+=GOOS=$(PI_OS) GOARCH=$(PI_ARCH) GOARM=$(PI_ARM)
endif

# General sources
SRCS=$(shell find . -name "*.go")

# Generated protobufs
PROTO_DIR=protos
PROTO_SRCS=flipdot/flipdot.pb.go flipapps/flipapps.pb.go
PROTO_BUFS=$(subst .pg.go,.proto,$(PROTO_SRCS))
MOCKS=$(subst .go,.mock.go,$(PROTO_SRCS)) flipdot/flipdot.mock.go

# Generated mocks
MOCKED_CLASS=FlipdotClient
MOCK_DIR=mock_flipdot
MOCK_FILE=$(MOCK_DIR)/flipdot.go

install: BUILD_CMD=go install
install: BUILD_ARGS=
install: build

build: protobufs
	$(ENVS) $(BUILD_CMD) $(BUILD_ARGS) $(MAIN_FILE)

protobufs: $(PROTO_SRCS)

test: $(MOCKS) $(PROTO_SRCS)
	go test ./...

%.pb.go: %.proto
	@echo Generating: $<
	protoc $(addprefix -I ,$(dir $(PROTO_BUFS))) --go_out=plugins=grpc:../../.. $<

%.mock.go: %.go
	mockgen -source $< -package $(lastword $(subst /, ,$(dir $<))) > $@

format:
	gofmt -w .

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_SRCS) $(MOCKS)

.PHONY: clean format test
