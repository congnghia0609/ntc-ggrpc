#!/bin/bash
# Author:       nghiatc
# Email:        congnghia0609@gmail.com

#source /etc/profile

echo "Install library dependencies..."

go get -u github.com/tools/godep
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
go get -u google.golang.org/grpc

echo "Install dependencies complete..."