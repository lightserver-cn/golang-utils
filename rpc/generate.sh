#!/bin/bash
echo "generate start"

# shellcheck disable=SC2164
cd protos

protoc --go_out=../pbs/ common.proto
protoc --micro_out=../pbs/ --go_out=../pbs/ common.proto
protoc-go-inject-tag -input=../pbs/common.pb.go

cd ../

echo "generate end"