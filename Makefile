MIGRATOR=go run main.go migrate

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

tidy: ## Tidy up dependencies, format code, and run vet
	go mod tidy
	go fmt ./...
	go vet ./...

dev: ## Run the API server in development mode
	air s

.PHONY: test
test: ## Run all tests
	go mod tidy
	go test -v --race

test_coverage: ## Run all tests with coverage
	go mod tidy
	go test -v -coverprofile=coverage.out --race

install: ## Install app to system
	go install

lint: ## Lint code
	golangci-lint run cmd/url_shortener

build: ## Build app
	go build -ldflags="-X code.commitHash=$(git rev-parse HEAD)" -o bin/url_shortener ./cmd/url_shortener

db-migrate: ## Run database migrations
	$(MIGRATOR) up

db-rollback: ## Rollback database migrations
	$(MIGRATOR) down

db-reset: ## Reset database to initial state
	$(MIGRATOR) reset

db-status: ## Show database migration status
	$(MIGRATOR) status

db-generate: ## Generate database code using sqlc
	docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc generate

docker-build:
	docker build -t url_shortener \
		--label "org.opencontainers.image.source=https://github.com/darkartx/go-project-278" \
		--label "org.opencontainers.image.description=Url shortener image" \
		.