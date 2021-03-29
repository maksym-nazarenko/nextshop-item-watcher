.PHONY: docker lint test

DOCKER_IMAGE ?= next-watcher
DOCKER_TAG ?= local
test:
	@go test ./...

lint:
	docker run --rm -v $$(git rev-parse --show-toplevel):/app:ro -w /app golangci/golangci-lint:v1.35.2 golangci-lint run -v

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f docker/Dockerfile .

run:
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.dev.yml up -d mongo

stop:
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.dev.yml stop mongo
