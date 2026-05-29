.PHONY: test race cover bench lint fmt vet check

# Run the test suite.
test:
	go test ./...

# Run tests with the race detector.
race:
	go test -race ./...

# Produce a coverage summary.
cover:
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Run benchmarks.
bench:
	go test -bench=. -benchmem -run=^$$ ./...

# Lint with golangci-lint.
lint:
	golangci-lint run ./...

# Format all Go files in place.
fmt:
	gofmt -w .

# go vet.
vet:
	go vet ./...

# Everything CI runs.
check: fmt vet test lint
