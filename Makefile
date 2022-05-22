BINARY_NAME=rated-cli

all: build

setup:
	go mod download

build:
	go build -v -o bin/$(BINARY_NAME) main.go


clean:
	rm -rf bin/*

.PHONY: all setup build clean
