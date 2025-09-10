.PHONY: all clean test install-deps get-deps generate generate-note-api install-golangci-lint lint lint-feature
all: lint

LOCAL_BIN:=$(CURDIR)/bin
install-deps:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
generate: install-deps
	make generate-note-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	api/note_v1/note.proto

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml

lint-feature:
	$(LOCAL_BIN)/golangci-lint run --config .golangci.yaml --new-from-rev dev