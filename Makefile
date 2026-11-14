BINARY_NAME = notse
MAIN_PATH   = ./cmd/notse
BUILD_DIR   = build

.PHONY: help all clean build install run

help:
	@echo "Available targets:"
	@echo " make help    - Show available targets"
	@echo " make all     - Run the full workflow"
	@echo " make clean   - Remove build artifacts"
	@echo " make build   - Build the binary"
	@echo " make install - Install the binary"
	@echo " make run     - Run the application"

all:
	@$(MAKE) --no-print-directory clean
	@$(MAKE) --no-print-directory build
	@$(MAKE) --no-print-directory install

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)/

build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

install:
	@echo "Installing binary..."
	@go install $(MAIN_PATH)

run:
	@echo "Running application..."
	@go run $(MAIN_PATH)
