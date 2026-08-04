.PHONY: build run clean install

BINARY_NAME=notse
MAIN_PATH=./cmd/notse

build:
	mkdir -p build
	go build -o build/$(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

clean:
	rm -rf build/

install:
	go install $(MAIN_PATH)
