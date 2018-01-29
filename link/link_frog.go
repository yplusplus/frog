// Copy this file to protoc-gen-go director, rebuild and reinstall protoc-gen-go.
// This file is to tell golang compiler to link frog plugin in protoc-gen-go.

package main

import _ "github.com/yplusplus/frog/plugin"
