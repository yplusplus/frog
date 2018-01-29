#!/bin/bash
protoc --go_out=plugins=frog,import_path=main,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. --proto_path=/path/to/protobuf/include:. *.proto
