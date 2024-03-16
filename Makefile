

## help: Display help
.PHONY: help
help: Makefile
	@echo "Usage:  make COMMAND"
	@echo
	@echo "Commands:"
	@sed -n 's/^##//p' $< | column -ts ':' |  sed -e 's/^/ /'

## tools: Install github.com/a-h/templ/cmd/templ@latest
.PHONY: tools
tools:
	go install github.com/a-h/templ/cmd/templ@latest

## get-deps: Download go dependencies
.PHONY: get-deps
get-deps:
	go mod download

## generate: Compile templ files
.PHONY: generate
generate:
	~/go/bin/templ generate

## build: Compile templ files and build application
.PHONY: build
build: get-deps generate
	CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o 'bin/app' ./cmd/app

## start: Build and start application
.PHONY: start
start: get-deps generate
	go run ./cmd/app
