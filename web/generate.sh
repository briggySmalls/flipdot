#!/bin/bash
set -e

# Establish some variables
GEN_OUTPUT_DIR=./src/generated
PROTOC_GEN_TS_PATH=$(yarn bin protoc-gen-ts)

# Copy protobufs to local
cp ../protos/flipdot.proto $GEN_OUTPUT_DIR/flipdot.proto
cp ../protos/flipapps.proto $GEN_OUTPUT_DIR/flipapps.proto

# Prepare output directory
mkdir -p $GEN_OUTPUT_DIR

# Build
yarn grpc_tools_node_protoc \
    -I $GEN_OUTPUT_DIR \
    --plugin="protoc-gen-ts=$PROTOC_GEN_TS_PATH" \
    --js_out=import_style=commonjs,binary:$GEN_OUTPUT_DIR \
    --ts_out=service=true:$GEN_OUTPUT_DIR \
    $GEN_OUTPUT_DIR/flipdot.proto $GEN_OUTPUT_DIR/flipapps.proto

# Remove protos
rm $GEN_OUTPUT_DIR/flipdot.proto $GEN_OUTPUT_DIR/flipapps.proto