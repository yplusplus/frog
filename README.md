# frog
Frog generates golang Google protobuf RPC interface. 
With the generated rpc interface, it is easy and flexibly to implement your own RPC system.

## Installation
```
+ go get -u github.com/yplusplus/frog
+ copy github.com/yplusplus/frog/link/link_frog.go to github.com/golang/protobuf/protoc-gen-go
+ re-build and re-install protoc-gen-go
```

And then, you can use following command to generate rpc code:
`protoc --go_out=plugins=frog:. *.proto`

## Examples
```
+ cd github.com/yplusplus/frog/example
+ go build
+ ./example
```
Check [example/](example/) for more detail.

## Other
Welcome to contribute

