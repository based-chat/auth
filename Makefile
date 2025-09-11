.PHONY: all install-deps generate generate-user-api install-golangci-lint lint lint-feature
all: lint

LOCAL_BIN:=$(CURDIR)/bin
install-deps:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
generate: install-deps
	make generate-user-api

generate-user-api:
	mkdir -p pkg/user/v1
	@if ! command -v protoc >/dev/null 2>&1 ; then echo "protoc not found в PATH"; exit 1; fi
	protoc \
	--proto_path api/user/v1 \
	--go_out=pkg/user/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/user/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	api/user/v1/user.proto

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

lint: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml

lint-feature: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run --config .golangci.yaml --new-from-rev dev