# frog
Frog generates golang Google protobuf RPC interface. 
With the generated rpc interface, it is easy and flexibly to implement your own RPC system.

## Installation
```
+ go get -u github.com/yplusplus/frog
+ copy github.com/yplusplus/frog/link/link_frog.go to github.com/golang/protobuf/protoc-gen-go
+ re-build and re-install protoc-gen-go
```

## Examples
```
+ cd github.com/yplusplus/frog/example
+ modify gen_proto.sh `--proto_path` to your protobuf include director
+ go generate
+ go build
+ ./example
```
Check [example/](example/) for more detail.

## Other
Welcome to contribute

