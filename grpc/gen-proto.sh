#!/bin/sh
path=$1
protoc --proto_path=$path --go_out=plugins=grpc:. $path/*.proto
