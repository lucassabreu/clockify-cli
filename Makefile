.PHONY: all
all: install-deps

export GO111MODULE=on

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## clean all buildable files
	rm -rf dist

install-deps: ## install golang dependencies
	go mod download

build: dist

dist: install-deps dist/darwin dist/linux dist/windows ## build all cli versions (default)

dist/darwin:
	mkdir -p dist/darwin
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/clockify-cli

dist/linux:
	mkdir -p dist/linux
	GOOS=linux GOARCH=amd64 go build -o dist/linux/clockify-cli

dist/windows:
	mkdir -p dist/windows
	GOOS=windows GOARCH=amd64 go build -o dist/windows/clockify-cli

go-install: ## install dev version
	go install

goreleaser-test:
	go install github.com/goreleaser/goreleaser
	goreleaser --snapshot --skip-publish --rm-dist
	go mod tidy

