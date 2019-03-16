OUTPUT=dist/main.js

FLIPDOT_DIR=./flipapp/flipdot
FLIPAPP_DIR=./flipapp/flipapps
FLIPAPP_PROTO=$(FLIPAPP_DIR)/flipapps.proto

PROTO_BUFS=flipdot.proto flipapps.proto
PROTO_SRCS_BEFORE_MOVE=$(patsubst %.proto,%_pb.js,$(PROTO_BUFS)) $(patsubst %.proto,%_grpc_web_pb.js,$(PROTO_BUFS))
PROTO_SRCS=$(notdir $(PROTO_SRCS_BEFORE_MOVE))
PROTO_DIR=../protos

app: $(OUTPUT)

proto: $(PROTO_SRCS)

$(OUTPUT): $(PROTO_SRCS)
	npx webpack client.js

%_pb.js %_grpc_web_pb.js: $(PROTO_BUFS)
	protoc \
		$(addprefix -I ,$(dir $(PROTO_BUFS))) \
		--js_out=import_style=commonjs:. \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:. \
		$^

%.proto:
	cp $(PROTO_DIR)/$(notdir $@) $@

clean:
	rm -f $(PROTO_SRCS)

# Helper to debug variables
print-%:
	@echo $*=$($*)

.PHONY: clean
