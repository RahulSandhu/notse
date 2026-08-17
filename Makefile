.PHONY: all clean build install run

BINARY_NAME=notse
MAIN_PATH=./cmd/notse

all: clean build install

clean:
	rm -rf build/

build:
	mkdir -p build
	go build -o build/$(BINARY_NAME) $(MAIN_PATH)

install:
	go install $(MAIN_PATH)

run:
	go run $(MAIN_PATH)
