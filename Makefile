.PHONY: all install-deps generate generate-user-api install-golangci-lint lint lint-feature clean test
all: clean generate build lint check-coverage  

LOCAL_BIN?=$(CURDIR)/.bin
PROTOC ?= protoc
BUILD_DIR ?= $(CURDIR)/bin

test: 
	go test -v -race ./...

clean:
	rm -rf $(LOCAL_BIN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	@rmdir pkg/user/v1 2>/dev/null || true

install-deps:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.9
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
generate: install-deps
	$(MAKE) generate-user-api

generate-user-api: install-deps
	mkdir -p pkg/user/v1
	@if ! command -v $(PROTOC) >/dev/null 2>&1 ; then \
		echo "Error: $(PROTOC) not found in PATH"; \
		echo "Please install protoc: https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	fi
	$(PROTOC) \
	--proto_path api/user/v1 \
	--go_out=pkg/user/v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go-grpc_out=pkg/user/v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	api/user/v1/user.proto

install-golangci-lint:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

lint: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.yaml

lint-fix: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run --fix ./... --config .golangci.yaml

lint-feature: install-golangci-lint
	$(LOCAL_BIN)/golangci-lint run --config .golangci.yaml --new-from-rev dev

.PHONY: install-go-test-coverage
install-go-test-coverage:
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./coverage.out  -covermode=atomic -coverpkg=./...
	$(LOCAL_BIN)/go-test-coverage --config=./.testcoverage.yml


build: generate build-server build-client

build-server: 
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/server cmd/grpc-server/main.go

build-client:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/client cmd/grpc-client/main.go

run-server: build-server
	$(BUILD_DIR)/server

run-client: build-client
	$(BUILD_DIR)/client