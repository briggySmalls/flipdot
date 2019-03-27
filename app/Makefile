# Build configuration
BIN_DIR?=bin
GO_CMD=go build
EXE_FILENAME=flipapp
LOCAL_EXE=$(BIN_DIR)/$(EXE_FILENAME)
GLOBAL_EXE=$(GOPATH)/bin/$(EXE_FILENAME)
BUILD_ARGS?=
DOCKER_TAG?=latest

# Cross-compilation
PI_OS=linux
PI_ARCH=arm
PI_ARM=7
IS_PI?=
ifdef IS_PI
	BUILD_ARGS+= -a
	ENVS+=GOOS=$(PI_OS) GOARCH=$(PI_ARCH) GOARM=$(PI_ARM)
endif

# Generated protobufs
PROTO_DIR=../protos
PROTO_SRCS=./flipdot/flipdot.pb.go ./flipapps/flipapps.pb.go
PROTO_BUFS=$(addprefix $(PROTO_DIR)/,$(notdir $(PROTO_SRCS:.pb.go=.proto)))
MOCK_SRCS=$(subst .go,.mock.go,$(PROTO_SRCS)) ./flipdot/flipdot.mock.go ./flipapps/button.mock.go ./flipapps/imager.mock.go

# Generated mocks
MOCKED_CLASS=FlipdotClient
MOCK_DIR=mock_flipdot
MOCK_FILE=$(MOCK_DIR)/flipdot.go

# All sources
TEST_SRCS=$(shell find . -name "*_test.go") $(MOCK_SRCS)
STD_SRCS=$(filter-out $(TEST_SRCS) $(PROTO_SRCS) $(MOCK_SRCS), $(shell find . -name "*.go"))
PROGRAM_SRCS=$(STD_SRCS) $(PROTO_SRCS)

## BUILD TARGETS

# Build to a local build directory
build: BUILD_ARGS+=-o $(LOCAL_EXE)
build: $(LOCAL_EXE)

# Build and install to GOPATH
install: $(PROGRAM_SRCS)
	$(ENVS) go install $(BUILD_ARGS) .

proto: $(PROTO_SRCS)

test: $(PROGRAM_SRCS) $(TEST_SRCS)
	go test ./...

format:
	gofmt -w .

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_SRCS) $(MOCK_SRCS)

docker:
	cp .dockerignore ..
	docker build -t briggysmalls/flipapp:$(DOCKER_TAG) -f Dockerfile ..
	rm ../.dockerignore

## PATTERN/SUB-RULES

$(LOCAL_EXE): $(PROGRAM_SRCS)
	$(ENVS) go build $(BUILD_ARGS) ./main.go

$(PROTO_SRCS): $(PROTO_BUFS)
	@echo Generating: $@
	for proto in $^; do \
		protoc $(addprefix -I , $(PROTO_DIR)) --go_out=plugins=grpc:$(GOPATH)/src $$proto; \
	done

%.mock.go: %.go
	mockgen -source $< -package $(lastword $(subst /, ,$(dir $<))) > $@

# Helper to debug variables
print-%:
	@echo $*=$($*)

.PHONY: clean format test docker
