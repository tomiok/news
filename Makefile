PACKAGES_PATH = $(shell go list -f '{{ .Dir }}' ./...)

.PHONY: all
all: require tidy goimports staticcheck test

.PHONY: require
require:
	@type "goimports" > /dev/null 2>&1 \
		|| (echo 'goimports not found: to install it, run "go install golang.org/x/tools/cmd/goimports@latest"'; exit 1)
	@type "staticcheck" > /dev/null 2>&1 \
		|| (echo 'staticcheck not found: to install it, run "go install honnef.co/go/tools/cmd/staticcheck@latest"'; exit 1)

.PHONY: tidy
tidy:
	@echo "=> Executing go mod tidy"
	@go mod tidy

.PHONY: goimports
goimports:
	@echo "=> Executing goimports"
	@goimports -w $(PACKAGES_PATH)


.PHONY: staticcheck
staticcheck:
	@echo "=> Executing staticcheck"
	@staticcheck ./...

.PHONY: test
test:
	@echo "=> Running tests"
	@go test ./... -count=1 -race -shuffle=on

.PHONY: test-cover
test-cover:
	@echo "=> Running tests and generating report"
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1 -race -shuffle=on
	@go tool cover -html=/tmp/coverage.out

.PHONY: run
run:
	@go run cmd/api/*.go

.PHONY: docker
docker:
	@docker compose up -d --no-deps --build

.PHONY: db
db:
	@docker compose up db