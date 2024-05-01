.PHONY: all
all: fmt tidy test

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -v ./... -bench=. -run=^$ -benchmem

.PHONY: clean
clean:
	@go clean && rm -rf ./bin/*
