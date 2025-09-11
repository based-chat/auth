//go:build tools

// Package tools фиксирует версии генераторов кода через blank-imports.
package tools

import (
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
