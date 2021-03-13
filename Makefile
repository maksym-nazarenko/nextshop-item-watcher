.PHONY: test lint

test:
	@go test ./...

lint:
	docker run --rm -v $$(git rev-parse --show-toplevel):/app:ro -w /app golangci/golangci-lint:v1.35.2 golangci-lint run -v
