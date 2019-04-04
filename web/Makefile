GEN_OUTPUT_DIR=src/generated

FLIPDOT_DIR=./flipapp/flipdot
FLIPAPP_DIR=./flipapp/flipapps
FLIPAPP_PROTO=$(FLIPAPP_DIR)/flipapps.proto

PROTO_BUFS=flipdot.proto flipapps.proto
PROTO_SRCS_BEFORE_MOVE=$(patsubst %.proto,%_pb.js,$(PROTO_BUFS)) $(patsubst %.proto,%_grpc_web_pb.js,$(PROTO_BUFS))
PROTO_SRCS=$(addprefix $(GEN_OUTPUT_DIR)/,$(notdir $(PROTO_SRCS_BEFORE_MOVE)))
PROTO_DIR=../protos

proto: $(PROTO_SRCS)

%_pb.js %_grpc_web_pb.js %_pb.d.ts: $(PROTO_BUFS)
	mkdir -p $(GEN_OUTPUT_DIR)
	protoc \
		$(addprefix -I ,$(dir $(PROTO_BUFS))) \
		--js_out=import_style=commonjs:$(GEN_OUTPUT_DIR) \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:$(GEN_OUTPUT_DIR) \
		$^

%.proto:
	cp $(PROTO_DIR)/$(notdir $@) $@

clean:
	rm -f $(PROTO_SRCS) $(GEN_OUTPUT_DIR)

# Helper to debug variables
print-%:
	@echo $*=$($*)

.PHONY: clean
