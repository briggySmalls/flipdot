#!/bin/bash
set -e

# Get current location
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Establish some variables
GEN_OUTPUT_DIR=$DIR/src/generated
PROTOC_GEN_TS_PATH=$(yarn bin protoc-gen-ts)

# Copy protobufs to local
mkdir -p $GEN_OUTPUT_DIR
cp $DIR/../protos/*.proto $GEN_OUTPUT_DIR

# Build
yarn grpc_tools_node_protoc \
    -I $GEN_OUTPUT_DIR \
    --plugin="protoc-gen-ts=$PROTOC_GEN_TS_PATH" \
    --js_out=import_style=commonjs,binary:$GEN_OUTPUT_DIR \
    --ts_out=service=true:$GEN_OUTPUT_DIR \
    $GEN_OUTPUT_DIR/flipdot.proto $GEN_OUTPUT_DIR/flipapps.proto

# Remove protos
rm $GEN_OUTPUT_DIR/*.proto