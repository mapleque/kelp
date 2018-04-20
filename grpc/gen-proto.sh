#!/bin/sh
protoc --proto_path=proto/ --go_out=plugins=grpc:. proto/*.proto
