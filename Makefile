SHELL := /bin/bash

PROJECT_NAME=openvpn-warder
BINARY_NAME=openvpn-warder
PACKAGE_NAME=github.com/qingstor/openvpn-warder

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  build          to build all"
	@echo "  clean          to clean the built file"

.PHONY: build

format:
	@echo "gofmt -w ."
	@gofmt -w .
	@echo "ok"

tidy:
	@go mod tidy
	@go mod verify
	@echo "ok"

build: format tidy
	@mkdir -p ./bin
	@echo "Building ${BINARY_NAME}..."
	@GOOS=${OS} GOARCH=amd64 go build -o ./bin/${BINARY_NAME}
	@echo "Done"

.PHONY: clean
clean:
	@echo "Clean the built file"
	rm -rf ./bin
	@echo "Done"
