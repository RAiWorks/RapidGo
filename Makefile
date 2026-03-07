.PHONY: build run test clean fmt vet lint

# Build the application binary
build:
	go build -o bin/rapidgo ./cmd/...

# Run the application
run:
	go run ./cmd/...

# Run all tests
test:
	go test ./... -v

# Remove build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run staticcheck (requires: go install honnef.co/go/tools/cmd/staticcheck@latest)
lint:
	staticcheck ./...
