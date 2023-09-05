#!/usr/bin/env bash

protoc -I=${GOPATH}/src:. --gogo_out=plugins=grpc,paths=source_relative:. *.proto
