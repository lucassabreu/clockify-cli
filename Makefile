.PHONY: all
all: install-deps

export GO111MODULE=on

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## clean all buildable files
	rm -rf build

install-deps: ## install golang dependencies
	go mod download

build: install-deps build/darwin build/linux build/windows ## build all cli versions (default)

build/darwin:
	mkdir -p build/darwin
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/clockify-cli

build/linux:
	mkdir -p build/linux
	GOOS=linux GOARCH=amd64 go build -o build/linux/clockify-cli

build/windows:
	mkdir -p build/windows
	GOOS=windows GOARCH=amd64 go build -o build/windows/clockify-cli
