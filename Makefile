IMAGE       := slots-backend:latest
BINARY      := backend
PROTO_TOOL  := buf
GO_TEST_OPTS := -race -cover ./...

.PHONY: help build run proto docker-build up down logs clean test

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  help           Show this help message"
	@echo "  build          Compile Go binary (./bin/$(BINARY))"
	@echo "  run            Run the service locally (go run ./cmd/server)"
	@echo "  test           Run unit + integration tests (go test $(GO_TEST_OPTS))"
	@echo "  proto          Regenerate gRPC code via $(PROTO_TOOL)"
	@echo "  docker-build   Build Docker image ($(IMAGE))"
	@echo "  up             Start services via docker-compose (detached)"
	@echo "  down           Stop and remove docker-compose services"
	@echo "  logs           Tail backend logs (last 50 lines)"
	@echo "  clean          Remove build artifacts (./bin)"

build:
	go build -o bin/$(BINARY) ./cmd/server

run:
	go run ./cmd/server

test:
	go test $(GO_TEST_OPTS)

proto:
	$(PROTO_TOOL) generate

docker-build:
	docker build -t $(IMAGE) .

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f --tail=50 backend

clean:
	rm -rf bin
