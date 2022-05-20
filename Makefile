BINARY_NAME=rated

all: build

setup:
	go mod download

build:
	go build -v -o bin/$(BINARY_NAME) ./cmd/main.go


clean:
	rm -rf bin/*

.PHONY: all setup build clean
