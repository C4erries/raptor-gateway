.PHONY: build run start
build:
		go build -v ./cmd/raptor-gateway
run:
	    go run ./cmd/raptor-gateway -config_path=./config/local.yaml
start: build run
		
		
.DEFAULT_GOAL := build