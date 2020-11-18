NAME=webhook-consumer
VERSION=dev
OS ?= linux
PROJECT_PATH ?= github.com/stone-co/webhook-consumer
PKG ?= github.com/stone-co/webhook-consumer/cmd
REGISTRY ?= stone-co
TERM=xterm-256color
CLICOLOR_FORCE=true
RICHGO_FORCE_COLOR=1

.PHONY: setup
setup:
	@echo "==> Setup: Getting tools"
	go mod tidy
	GO111MODULE=on go install \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/kyoh86/richgo

.PHONY: test
test:
	@echo "==> Running Tests"
	go test -v ./...

.PHONY: compile
compile: clean
	@echo "==> Go Building WebHookConsumer"
	@env GOOS=${OS} GOARCH=amd64 go build -v -o build/${NAME} ${PKG}/

.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=${OS} go clean -i -x ./...
	@rm -f build/${NAME}

.PHONY: metalint
metalint:
	@echo "==> Running Linters"
	go test -i ./...
	golangci-lint run -c ./.golangci.yml ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests"
	@richgo test -failfast -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
