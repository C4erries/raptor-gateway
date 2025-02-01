.PHONY: build
build:
		go build -v ./cmd/raptor-gateway

.DEFAULT_GOAL := build