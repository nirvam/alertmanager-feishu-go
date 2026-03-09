BINARY_NAME=alertmanager-feishu-go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -s -w"

.PHONY: all build test lint clean run

all: lint test build

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/${BINARY_NAME} main.go

test:
	go test -v ./...

lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

run:
	go run main.go serve

clean:
	rm -rf bin/
	rm -rf vendor/

install: build
	install -Dm755 bin/${BINARY_NAME} /usr/bin/${BINARY_NAME}
