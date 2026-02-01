# Build variables
BINARY_NAME=masstdb
VERSION?=0.1.0
BUILD_DIR=dist

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test install release

# Default target
all: build

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
test:
	go test -v ./...

# Install locally
install:
	go install $(LDFLAGS) .

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

# Build for all platforms
release: clean
	mkdir -p $(BUILD_DIR)
	# macOS (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# macOS (Intel)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# Linux (amd64)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# Linux (arm64)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Binaries built in $(BUILD_DIR)/"
	@ls -la $(BUILD_DIR)/
