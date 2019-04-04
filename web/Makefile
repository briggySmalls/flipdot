GEN_OUTPUT_DIR=src/generated

FLIPDOT_DIR=./flipapp/flipdot
FLIPAPP_DIR=./flipapp/flipapps
FLIPAPP_PROTO=$(FLIPAPP_DIR)/flipapps.proto

PROTO_BUFS=$(addprefix $(GEN_OUTPUT_DIR)/,flipdot.proto flipapps.proto)
PROTO_SRCS_BEFORE_MOVE=$(patsubst %.proto,%_pb.js,$(PROTO_BUFS)) $(patsubst %.proto,%_grpc_web_pb.js,$(PROTO_BUFS))
PROTO_SRCS=$(notdir $(PROTO_SRCS_BEFORE_MOVE))
PROTO_DIR=../protos

# Path to this plugin
PROTOC_GEN_TS_PATH=./node_modules/.bin/protoc-gen-ts

proto: $(PROTO_SRCS)

%_pb.js %_grpc_web_pb.js %_pb.d.ts: $(PROTO_BUFS)
	protoc \
		-I $(GEN_OUTPUT_DIR) \
		--plugin="protoc-gen-ts=$(PROTOC_GEN_TS_PATH)" \
		--js_out=import_style=commonjs,binary:$(GEN_OUTPUT_DIR) \
		--ts_out=service=true:$(GEN_OUTPUT_DIR) \
		$^

%.proto:
	mkdir -p $(GEN_OUTPUT_DIR)
	cp $(PROTO_DIR)/$(notdir $@) $@

clean:
	rm -rf $(GEN_OUTPUT_DIR)

# Helper to debug variables
print-%:
	@echo $*=$($*)

.PHONY: clean
