.PHONY: build run clean test install deps

BINARY_NAME=notse
MAIN_PATH=./cmd/notse

build:
	go build -o ${BINARY_NAME} ${MAIN_PATH}

run:
	go run ${MAIN_PATH}

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test -v ./...

install:
	go install ${MAIN_PATH}

deps:
	./scripts/install-deps.sh

help:
	@echo "Available commands:"
	@echo "  make deps    - Install dependencies"
	@echo "  make build   - Build the binary"
	@echo "  make run     - Run the app"
	@echo "  make install - Install to GOPATH/bin"
	@echo "  make clean   - Remove built files"
	@echo "  make test    - Run tests"
