.PHONY: build run start
build:
		go build -v ./cmd/raptor-gateway
		sleep 5
run:
		go run ./cmd/raptor-gateway
start: build run
		
		
.DEFAULT_GOAL := build