
# -------------------------------------------------------------------------------------------------
# main
# -------------------------------------------------------------------------------------------------

all: help

## build: Compile templ files and build application
.PHONY: build
build: get-deps generate
	CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o 'bin/app' ./cmd/app

## start: Build and start application
.PHONY: start
start: get-deps generate
	go run ./cmd/app

## build-docker: Build Docker container image with this app
.PHONY: build-docker
build-docker:
	docker build -t $(shell basename $(PWD)):latest --no-cache -f Dockerfile .

## run-docker: Run Docker container image with this app
.PHONY: run-docker
run-docker:
	docker run --rm -it -p 8080:8080 $(shell basename $(PWD)):latest

# -------------------------------------------------------------------------------------------------
# testing
# -------------------------------------------------------------------------------------------------

## test: Run unit tests
.PHONY: test
test: check-go
	@go test ./...

# -------------------------------------------------------------------------------------------------
# tools && shared
# -------------------------------------------------------------------------------------------------

## tidy: Removes unused dependencies and adds missing ones
.PHONY: tidy
tidy: check-go
	go mod tidy

## get-deps: Download go dependencies
.PHONY: get-deps
get-deps: check-go
	go mod download

## tools: Install github.com/a-h/templ/cmd/templ@latest
.PHONY: tools
tools: check-go
	go install github.com/a-h/templ/cmd/templ@latest

## get-air: Install live reload server github.com/cosmtrek/air@latest
.PHONY: get-air
get-air: check-go
	go install github.com/cosmtrek/air@latest

## generate: Compile templ files
.PHONY: generate
generate:
	~/go/bin/templ generate

## air: Build and start application in live reload mode via air
.PHONY: air
air: get-deps generate
	air

## check-go: Check that Go is installed
.PHONY: check-go
check-go:
	@command -v go &> /dev/null || (echo "Please install GoLang" && false)

## help: Display help
.PHONY: help
help: Makefile
	@echo "Usage:  make COMMAND"
	@echo
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -ts ':' |  sed -e 's/^/ /'
