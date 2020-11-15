

build:
	go build ./cmd/secretize

test:
	@go test -v -cover ./...
	@echo "All tests passed"