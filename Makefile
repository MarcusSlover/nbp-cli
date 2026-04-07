# Build script for NBP CLI tool

BINARY_NAME=nbp

# Local development build
dev:
	go build -o $(BINARY_NAME) main.go

# Run tests
test:
	go test ./...

# Standard installation using go install
install:
	go install

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Run goreleaser locally for verification (snapshot)
release-snapshot:
	goreleaser release --snapshot --clean

.PHONY: dev test install clean release-snapshot
